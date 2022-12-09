// Command xelf provides repl and other tools to work with xelf languages.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"xelf.org/cmd"
	"xelf.org/xelf/xps"
)

var dirFlag = flag.String("dir", ".", "")

func main() {
	flag.Usage = printUsage
	flag.Parse()
	log.SetFlags(0)
	ctx := cmd.DefaultCtx(*dirFlag, flag.Args())
	subcmd := ctx.Split()
	switch subcmd {
	case "version":
		fmt.Println("alpha-dev")
	case "bash.inc":
		printCompletion(ctx)
	default:
		c := cmd.All[subcmd]
		if c == nil {
			plugCmd, err := xps.PlugCmd(ctx, subcmd)
			if err != nil {
				log.Fatalf("loading plug %s: %v", subcmd, err)
			}
			c = &cmd.Cmd{Name: subcmd, Func: plugCmd}
		}
		if c == nil {
			log.Printf("unknown subcommand %q\n", subcmd)
			fmt.Print(usage)
			printPluginHelp(ctx)
			os.Exit(1)
		}
		err := c.Func(ctx)
		if err != nil {
			if redir, _ := err.(*xps.CmdRedir); redir != nil {
				c = cmd.All[redir.Cmd]
				if err = c.Func(ctx); err != nil {
					log.Fatalf("xelf %s→ %s %v", subcmd, redir.Cmd, err)
				}
			} else {
				log.Fatalf("xelf %s %v", subcmd, err)
			}
		}
	case "":
		printUsage()
	case "help":
		fmt.Print(usage)
		printPluginHelp(ctx)
		fmt.Println()
	}
}

func printUsage() {
	fmt.Print(usage)
	fmt.Println("For plug-in information run xelf help.")
}

const usage = `Usage: xelf [-dir=<path>] <command> [<args>]

Most commands do read from stdin and print to standard output and error.
Use pipes and redirects to use strings or read or write to files:

   xelf sel 'posts/title' < blog.xelf | xelf fmt > titles.json
   echo '{a:1}{a:2}' | xelf sel a
   xelf sel a <<< '{a:1}{a:2}'

The -dir flag specifies the base dir and defaults to the current directory.

Evaluation commands accept arguments or stdin as input
   run         Evaluates the input and prints errors or the result.
   test        Resolves the input and prints errors or the result type.
   repl        Starts a read-eval-print-loop to explore xelf. It uses the input
               as preamble and makes $XELF_PLUGINS available for import.

Development commands
   fmt         Formats the input and prints the result.
   fix         Fixes the input and prints the result.
   list        Searches for files and modules and prints the result.
   rebuild     Rebuilds all plugin sources it can find.

Literal commands work on streams of json or xelf values
   sel         Selects a path from the input and streams the result.
   mut         Applies a delta to the input and steams the result.
   json        Converts the input to JSON and streams the result.

Other commands
   version     Prints xelf version information.
   bash.inc    Prints xelf bash completion including plugin subcommands.
               The output can be saved to a file and sourced by bashrc.
   help        Displays this help message.

`

func printPluginHelp(ctx *xps.CmdCtx) {
	var header bool
	ms := ctx.Manifests()
	for _, m := range ms {
		cmds := m.Cmds()
		if len(cmds) == 0 {
			continue
		}
		if !header {
			fmt.Println("Plug-in commands")
			header = true
		}
		fmt.Println()
		var info string
		if fst := cmds[0]; fst.Key == "" {
			info = fst.Val.String()
			cmds = cmds[1:]
		} else {
			info = fmt.Sprintf("Subcommands for %s", m.Name)
		}
		fmt.Printf("   %-11s %s\n", m.Name, info)
		for _, sub := range cmds {
			if sub.Key == "" { // group heading
				fmt.Printf("   * %s\n", sub.Val)
			} else {
				fmt.Printf("   · %-9s %s\n", sub.Key, sub.Val)
			}
		}
	}
}

func printCompletion(ctx *xps.CmdCtx) {
	var subs, plugs, psub []string
	subs = append(subs, "version", "help")
	for sub := range cmd.All {
		subs = append(subs, sub)
	}
	for _, m := range ctx.Manifests() {
		cmds := m.Cmds()
		if len(cmds) == 0 {
			continue
		}
		psub = psub[:0]
		for _, kv := range cmds {
			if kv.Key != "" {
				psub = append(psub, kv.Key)
			}
		}
		subs = append(subs, m.Name)
		if len(psub) == 0 {
			continue
		}
		plugs = append(plugs, fmt.Sprintf(plugFmt, m.Name, strings.Join(psub, " ")))
	}
	fmt.Printf(complFmt, strings.Join(subs, " "), strings.Join(plugs, ""))
}

var complFmt = `_xelf_complete() {
	set -- $COMP_LINE
	shift; while [[ $1 == -* ]]; do shift; done
	local cur=${COMP_WORDS[COMP_CWORD]}%[2]s
	COMPREPLY=( $(compgen -W "%[1]s" -- $cur) )
}
complete -F _xelf_complete xelf
`

var plugFmt = `
	if grep -q '^\(%[1]s\)$' <<< $1; then
		COMPREPLY=( $(compgen -W "%[2]s" -- $cur) ); return
	fi`

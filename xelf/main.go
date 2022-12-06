// Command xelf provides repl and other tools to work with xelf languages.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"xelf.org/cmd"
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	subcmd := flag.Arg(0)
	switch subcmd {
	case "version":
		fmt.Println("alpha-dev")
	default:
		c := cmd.All[subcmd]
		if c == nil {
			log.Printf("unknown subcommand %q\n", subcmd)
			fmt.Print(usage)
			os.Exit(1)
		}
		err := c.Func(flag.Args()[1:])
		if err != nil {
			log.Fatalf("xelf %s %v", subcmd, err)
		}
	case "", "help":
		fmt.Print(usage)
	}
}

const usage = `usage: xelf <command> [<args>]

Most commands do read from standard input and print to standard output and error.
Use pipes and redirects to read or write to files:

   xelf sel 'posts/title' < blog.xelf | xelf fmt > titles.json

Evaluation commands accept arguments or stdin as input:

   run         Evaluates the input and prints errors to stderr or the result to stdout.
   test        Resolves the input and prints errors to stderr or the result type to stdout.
   repl        Starts a read-eval-print-loop to explore xelf with. Input acts as program init.
               The repl looks for $XELF_PLUGINS and makes them available for import.

Development commands

   fmt         Formats the input and prints the result.
   fix         Fixes and formats the input or argument file and prints the result.
   list        Searches for files and modules and prints the result.

Literal commands

   sel         Selects a path from input or an argument file and prints the result.
   mut         Applies a delta to the input or an argument file and prints the result.
   json        Converts the input or argument file to JSON and prints the result.


Other commands

   version     Prints xelf version information.
   help        Displays this help message.

`

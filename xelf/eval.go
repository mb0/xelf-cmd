package main

import (
	"fmt"
	"os"
	"strings"

	"xelf.org/cmd"
	"xelf.org/xelf/exp"
	"xelf.org/xelf/lit"
)

var _ = cmd.Add("run", func(args []string) error {
	return cmd.SafetyWrap(func() (err error) {
		var x exp.Exp
		if len(args) > 0 {
			x, err = exp.Parse(strings.Join(args, " "))
		} else {
			x, err = exp.Read(os.Stdin, "stdin")
		}
		if err != nil {
			return err
		}
		res, err := cmd.Prog().Run(x, &lit.Keyed{})
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	})
})
var _ = cmd.Add("repl", func(args []string) error {
	return cmd.SafetyWrap(func() error {
		r := cmd.NewRepl(cmd.ReplHistoryPath("xelf/repl.history"))
		defer r.Close()
		r.Run()
		return nil
	})
})

package main

import (
	"log"

	"xelf.org/cmd"
)

func eval(args []string) {
	res, err := cmd.Prog().RunStr(args[0], nil)
	log.Printf("%v %v", res, err)
}
func repl(args []string) error {
	r := cmd.NewRepl(cmd.ReplHistoryPath("xelf/repl.history"))
	defer r.Close()
	r.Run()
	return nil
}

package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	args := flag.Args()
	switch flag.Arg(0) {
	case "help":
		help()
	case "eval":
		eval(args[1:])
	case "repl":
		repl(args[1:])
	default:
		help()
	}
}

func help() {
	fmt.Println("xelf command")
}

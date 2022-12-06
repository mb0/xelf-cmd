package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"xelf.org/cmd"
	"xelf.org/xelf/ast"
	"xelf.org/xelf/cor"
	"xelf.org/xelf/lit"
)

var _ = cmd.Add("sel", func(dir string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("expects a selection path argument")
	}
	// read path from first arg
	p, err := cor.ParsePath(args[0])
	if err != nil {
		return fmt.Errorf("could not read the path: %v", err)
	}
	// optionally followed by path variables
	if n := p.CountVars(); n > 0 {
		if err = p.FillVars(args[1:]); err != nil {
			return err
		}
	}
	// apply the selection to a stream of values
	return stdinStream(func(val lit.Val) error {
		// select the path
		res, err := lit.SelectPath(val, p)
		if err != nil {
			return err
		}
		// and print results
		fmt.Println(res)
		return nil
	})
})

var _ = cmd.Add("mut", func(dir string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("expects a delta dict argument")
	}
	var delta lit.Delta
	err := lit.ParseInto(args[0], (*lit.Keyed)(&delta))
	if err != nil {
		return fmt.Errorf("could not read the delta dict: %v", err)
	}
	// apply the mutation to a stream of values
	return stdinStream(func(val lit.Val) error {
		mut := val.Mut()
		mut, err = lit.Apply(mut, delta)
		if err != nil {
			return fmt.Errorf("failed: %v", err)
		}
		fmt.Println(mut)
		return nil
	})
})

var _ = cmd.Add("json", func(dir string, args []string) error {
	// TODO args maybe to configure pretty printing?
	// print a stream of values as json
	return stdinStream(func(val lit.Val) error {
		buf, err := val.MarshalJSON()
		if err != nil {
			return fmt.Errorf("could not marshal literal: %v", err)
		}
		os.Stdout.Write(buf)
		fmt.Println()
		return nil
	})
})

func stdinStream(f func(lit.Val) error) error {
	// create a new lexer for stdin
	lex := ast.NewLexer(os.Stdin, "stdin")
	// apply the function to a stream of values
	for {
		// read an ast node from the lexer
		a, err := ast.Scan(lex)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("could not read a literal from stdin: %v", err)
		}
		// parse value
		val, err := lit.ParseVal(a)
		if err != nil {
			return err
		}
		// apply the function
		err = f(val)
		if err != nil {
			return fmt.Errorf("failed: %v", err)
		}
	}
	return nil
}

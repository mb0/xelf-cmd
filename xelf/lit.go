package main

import (
	"fmt"
	"os"

	"xelf.org/cmd"
	"xelf.org/xelf/cor"
	"xelf.org/xelf/lit"
)

var _ = cmd.Add("sel", func(args []string) error {
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
	// read the literal from stdin
	val, err := lit.Read(os.Stdin, "stdin")
	if err != nil {
		return fmt.Errorf("could not read a literal from stdin: %v", err)
	}
	// selection the path
	res, err := lit.SelectPath(val, p)
	if err != nil {
		return fmt.Errorf("failed: %v", err)
	}
	// and print results
	fmt.Println(res)
	return nil
})

var _ = cmd.Add("mut", func(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("expects a delta dict argument")
	}
	var delta lit.Delta
	err := lit.ParseInto(args[0], (*lit.Keyed)(&delta))
	if err != nil {
		return fmt.Errorf("could not read the delta dict: %v", err)
	}
	val, err := lit.Read(os.Stdin, "stdin")
	if err != nil {
		return fmt.Errorf("could not read a literal from stdin: %v", err)
	}
	mut := val.Mut()
	mut, err = lit.Apply(mut, delta)
	if err != nil {
		return fmt.Errorf("failed: %v", err)
	}
	fmt.Println(mut)
	return nil
})

var _ = cmd.Add("json", func(args []string) error {
	// TODO args maybe to configure pretty printing?
	val, err := lit.Read(os.Stdin, "stdin")
	if err != nil {
		return fmt.Errorf("could not read a literal from stdin: %v", err)
	}
	buf, err := val.MarshalJSON()
	if err != nil {
		return fmt.Errorf("could not marshal literal: %v", err)
	}
	os.Stdout.Write(buf)
	return nil
})

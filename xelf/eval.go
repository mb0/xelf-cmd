package main

import (
	"fmt"
	"os"
	"strings"

	"xelf.org/cmd"
	"xelf.org/xelf/exp"
	"xelf.org/xelf/lit"
	"xelf.org/xelf/typ"
	"xelf.org/xelf/xps"
)

var _ = cmd.Add("run", func(ctx *xps.CmdCtx) error {
	return cmd.SafetyWrap(func() (err error) {
		var x exp.Exp
		if len(ctx.Args) > 0 {
			x, err = exp.Parse(strings.Join(ctx.Args, " "))
		} else {
			x, err = exp.Read(os.Stdin, "stdin")
		}
		if err != nil {
			return err
		}
		res, err := cmd.DefaultProg(ctx).Run(x, &lit.Keyed{})
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	})
})
var _ = cmd.Add("test", func(ctx *xps.CmdCtx) error {
	return cmd.SafetyWrap(func() (err error) {
		var x exp.Exp
		if len(ctx.Args) > 0 {
			x, err = exp.Parse(strings.Join(ctx.Args, " "))
		} else {
			x, err = exp.Read(os.Stdin, "stdin")
		}
		if err != nil {
			return err
		}
		p := cmd.DefaultProg(ctx)
		p.Arg = &lit.Keyed{}
		res, err := p.Resl(p, x, typ.Void)
		if err != nil {
			return err
		}
		fmt.Println(res.Type())
		return nil
	})
})
var _ = cmd.Add("repl", func(ctx *xps.CmdCtx) error {
	return cmd.SafetyWrap(func() error {
		hist := &cmd.FileReplHist{Path: cmd.ReplHistoryPath("xelf/repl.history")}
		r := cmd.NewRepl(ctx, hist)
		defer r.Close()
		r.Run()
		return nil
	})
})

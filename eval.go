package cmd

import (
	"xelf.org/xelf/exp"
	"xelf.org/xelf/lib/extlib"
	"xelf.org/xelf/mod"
	"xelf.org/xelf/xps"
)

func DefaultProg(ctx *xps.CmdCtx) *exp.Prog {
	if ctx.Prog != nil {
		return ctx.Prog(ctx)
	}
	return exp.NewProg(DefaultEnv(ctx))
}

func DefaultEnv(ctx *xps.CmdCtx) exp.Env {
	ctx.Manifests()
	return mod.NewLoaderEnv(extlib.Std,
		&xps.ModLoader{Sys: mod.Registry, Plugs: &ctx.Plugs},
		mod.FileMods(ctx.Dir),
	)
}

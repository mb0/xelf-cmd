package cmd

import (
	"xelf.org/xelf/exp"
	"xelf.org/xelf/lib/extlib"
	"xelf.org/xelf/mod"
	"xelf.org/xelf/xps"
)

func ProgRoot() exp.Env {
	roots := xps.EnvRoots()
	mods := xps.NewMods(mod.Registry, xps.FindAll(roots))
	fmods := mod.FileMods()
	return mod.NewLoaderEnv(extlib.Std, mods, fmods)
}

func Prog() *exp.Prog {
	return exp.NewProg(ProgRoot())
}

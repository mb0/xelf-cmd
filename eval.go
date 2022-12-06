package cmd

import (
	"xelf.org/xelf/exp"
	"xelf.org/xelf/lib/extlib"
	"xelf.org/xelf/mod"
	"xelf.org/xelf/xps"
)

func Prog() *exp.Prog {
	roots := xps.EnvRoots()
	mods := xps.NewMods(mod.Registry, xps.FindAll(roots))
	fmods := mod.FileMods()
	env := mod.NewLoaderEnv(extlib.Std, mods, fmods)
	return exp.NewProg(env)
}

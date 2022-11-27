package cmd

import (
	"fmt"

	"xelf.org/xelf/exp"
	"xelf.org/xelf/lib/extlib"
	"xelf.org/xelf/mod"
	"xelf.org/xelf/xps"
)

func Prog() *exp.Prog {
	roots := xps.EnvRoots()
	plugs, err := xps.LoadAll(roots)
	if err != nil {
		fmt.Printf("failed to load plugins: %v\n", err)
	}
	fmt.Printf("found plugins %s in %s\n", plugs, roots)
	fmods := mod.FileMods()
	env := mod.NewLoaderEnv(extlib.Std, mod.Registry, fmods)
	return exp.NewProg(env)
}

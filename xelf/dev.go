package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"xelf.org/cmd"
	"xelf.org/xelf/xps"
)

var _ = cmd.Add("fmt", func(ctx *xps.CmdCtx) error {
	// TODO implements real formatter for now only pipe stdin to stdout
	_, err := io.Copy(os.Stdout, os.Stdin)
	return err
})
var _ = cmd.Add("fix", func(ctx *xps.CmdCtx) error {
	// TODO implements real formatter for now only pipe stdin to stdout
	_, err := io.Copy(os.Stdout, os.Stdin)
	return err
})
var _ = cmd.Add("list", func(ctx *xps.CmdCtx) error {
	// TODO use fsmods to discover all xelf files in dir
	pms := ctx.Manifests()
	if len(pms) == 0 {
		if roots := xps.EnvRoots(); len(roots) == 0 {
			fmt.Println("No XELF_PLUGINS path list set.")
			return nil
		} else {
			fmt.Printf("No plugin manifests found in: %s\n", strings.Join(roots, ", "))
		}
		return nil
	} else {
		fmt.Println("Found plugin manifests:")
		for _, pm := range pms {
			fmt.Printf("   %-11s (%s)\n", pm.Name, pm.Path)
			const dotfmt = "   · %-9s %s\n"
			if mods := pm.Mods(); len(mods) > 0 {
				fmt.Printf(dotfmt, "Modules:", strings.Join(mods, ", "))
			}
			cs := pm.Cmds()
			if len(cs) > 0 {
				cmds := make([]string, 0, len(cs))
				for _, c := range cs {
					if c.Key != "" {
						name := fmt.Sprintf("%s %s", pm.Name, c.Key)
						cmds = append(cmds, name)
					}
				}
				if len(cmds) == 0 {
					cmds = append(cmds, pm.Name)
				}
				fmt.Printf(dotfmt, "Commands:", strings.Join(cmds, ", "))
			}
		}
	}
	return nil
})
var _ = cmd.Add("rebuild", func(ctx *xps.CmdCtx) error {
	// TODO use args to filter plugins?
	pms := ctx.Manifests()
	fmt.Printf("Checking %d plugins…\n", len(pms))
	var errn int
	for _, pm := range pms {
		if !xps.HasSource(pm) {
			fmt.Printf("   · no source for %s\n", pm.Name)
			continue
		}
		fmt.Printf("   · building %s …\n", pm.Name)
		err := xps.Rebuild(pm)
		if err != nil {
			errn++
			fmt.Printf("   ! failed: %v\n", err)
			continue
		}
		fmt.Printf("     ok\n")
	}
	if errn != 0 {
		return fmt.Errorf("failed to rebuild some plugin sources")
	}
	return nil
})

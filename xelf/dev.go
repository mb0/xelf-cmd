package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"xelf.org/cmd"
	"xelf.org/xelf/xps"
)

var _ = cmd.Add("fmt", func(dir string, args []string) error {
	// TODO implements real formatter for now only pipe stdin to stdout
	_, err := io.Copy(os.Stdout, os.Stdin)
	return err
})
var _ = cmd.Add("fix", func(dir string, args []string) error {
	// TODO implements real formatter for now only pipe stdin to stdout
	_, err := io.Copy(os.Stdout, os.Stdin)
	return err
})
var _ = cmd.Add("list", func(dir string, args []string) error {
	// TODO use fsmods to discover all xelf files in dir
	roots := xps.EnvRoots()
	if len(roots) == 0 {
		fmt.Println("No XELF_PLUGINS path list set.")
		return nil
	}
	pms := xps.FindAll(roots)
	if len(pms) == 0 {
		fmt.Printf("No plugin manifests found in: %s\n", strings.Join(roots, ", "))
		return nil
	} else {
		fmt.Println("Found plugin manifests:")
		for _, pm := range pms {
			fmt.Printf("   %-11s (%s)\n", pm.Name, pm.Path)
			const dotfmt = "   · %-9s %s\n"
			if len(pm.Mods) > 0 {
				fmt.Printf(dotfmt, "Modules:", strings.Join(pm.Mods, ", "))
			}
			if len(pm.Cmds) > 0 {
				cmds := make([]string, 0, len(pm.Cmds))
				for _, c := range pm.Cmds {
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

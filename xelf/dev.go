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
	roots := xps.EnvRoots()
	if fmt.Print("XELF_PLUGINS "); len(roots) == 0 {
		fmt.Println("not set")
	} else {
		fmt.Println("set to", strings.Join(roots, ", "))
	}
	pms := xps.FindAll(xps.EnvRoots())
	if len(pms) == 0 {
		fmt.Println("No manifest files found")
	} else {
		fmt.Println("Manifest files:")
		for _, pm := range pms {
			fmt.Printf("\t%s provides: %s\n", pm.Path, strings.Join(pm.Mods, ", "))
		}
	}
	return nil
})

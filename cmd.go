// Package cmd provides helpers for other commands working with xelf.
// The actual command can be found in the sub directory at xelf.org/cmd/xelf.
package cmd

import "fmt"

// SafetyWrap wraps a function and provides panic recovery.
func SafetyWrap(f func() error) error {
	errc := make(chan error)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				errc <- fmt.Errorf("panic: %v", err)
			}
		}()
		errc <- f()
	}()
	return <-errc
}

type Cmd struct {
	Name string
	Func func(args []string) error
}

var All = make(map[string]*Cmd)

func Add(name string, f func([]string) error) *Cmd {
	c := &Cmd{name, f}
	All[name] = c
	return c
}

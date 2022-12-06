package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/peterh/liner"
	"xelf.org/xelf/bfr"
	"xelf.org/xelf/exp"
	"xelf.org/xelf/lit"
)

func ReplHistoryPath(rest string) string {
	path, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	return filepath.Join(path, rest)
}

type Repl struct {
	*liner.State
	Hist string
	Prog *exp.Prog
	Wrap func(exp.Env) exp.Env
}

func NewRepl(hist string) *Repl {
	lin := liner.NewLiner()
	lin.SetMultiLineMode(true)
	return &Repl{State: lin, Hist: hist}
}

func (r Repl) Run() {
	defer r.Close()
	r.readHistory()
	var raw []byte
	arg := &lit.Keyed{}
	p := r.Prog
	if p == nil {
		p = Prog()
	}
	for {
		prompt := "> "
		if len(raw) > 0 {
			prompt = "… "
		}
		got, err := r.Prompt(prompt)
		if err != nil {
			if err == io.EOF {
				r.writeHistory()
				fmt.Println()
				return
			}
			raw = raw[:0]
			log.Printf("unexpected error reading prompt: %v", err)
			continue
		}
		raw = append(raw, got...)
		el, err := exp.Read(bytes.NewReader(raw), "")
		if err != nil {
			if errors.Is(err, io.EOF) {
				raw = append(raw, '\n')
				continue
			}
			log.Printf("! %v\n\n", err)
			r.AppendHistory(string(raw))
			raw = raw[:0]
			continue
		}
		r.AppendHistory(bfr.String(el))
		raw = raw[:0]
		org := p.Root
		if r.Wrap != nil {
			p.Root = r.Wrap(org)
		}
		l, err := p.Run(el, arg)
		if r.Wrap != nil {
			p.Root = org
		}
		if err != nil {
			log.Printf("! %v\n\n", err)
			continue
		}
		fmt.Printf("= %s\n\n", bfr.String(l))
	}
}

func (r *Repl) readHistory() {
	if r.Hist == "" {
		return
	}
	f, err := os.Open(r.Hist)
	if err != nil {
		log.Printf("error reading repl history file %q: %v\n", r.Hist, err)
		return
	}
	defer f.Close()
	_, err = r.ReadHistory(f)
	if err != nil {
		log.Printf("error reading repl history file %q: %v\n", r.Hist, err)
	}
}

func (r *Repl) writeHistory() {
	if r.Hist == "" {
		return
	}
	dir := filepath.Dir(r.Hist)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Printf("error creating dir for repl history %q: %v\n", dir, err)
		return
	}
	f, err := os.Create(r.Hist)
	if err != nil {
		log.Printf("error creating file for repl history %q: %v\n", r.Hist, err)
		return
	}
	defer f.Close()
	_, err = r.WriteHistory(f)
	if err != nil {
		log.Printf("error writing repl history file %q: %v\n", r.Hist, err)
	}
}

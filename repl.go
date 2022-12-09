package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/peterh/liner"
	"xelf.org/xelf/bfr"
	"xelf.org/xelf/exp"
	"xelf.org/xelf/lit"
	"xelf.org/xelf/xps"
)

type Repl struct {
	*liner.State
	Hist ReplHistory
	Ctx  *xps.CmdCtx
}

func NewRepl(ctx *xps.CmdCtx, hist ReplHistory) *Repl {
	lin := liner.NewLiner()
	lin.SetMultiLineMode(true)
	return &Repl{State: lin, Ctx: ctx, Hist: hist}
}

func (r *Repl) Run() {
	defer r.Close()
	err := r.Hist.Read(r)
	if err != nil {
		log.Printf("error reading repl history: %v\n", err)
	}
	var raw []byte
	arg := &lit.Keyed{}
	p := DefaultProg(r.Ctx)
	for {
		prompt := "> "
		if len(raw) > 0 {
			prompt = "… "
		}
		got, err := r.Prompt(prompt)
		if err != nil {
			if err == io.EOF {
				herr := r.Hist.Write(r)
				if herr != nil {
					log.Printf("error writing repl history: %v\n", herr)
				}
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
		if r.Ctx.Wrap != nil {
			p.Root = r.Ctx.Wrap(r.Ctx, org)
		}
		l, err := p.Run(el, arg)
		if r.Ctx.Wrap != nil {
			p.Root = org
		}
		if err != nil {
			log.Printf("! %v\n\n", err)
			continue
		}
		fmt.Printf("= %s\n\n", bfr.String(l))
	}
}

package cmd

import (
	"os"
	"path/filepath"
)

func ReplHistoryPath(rest string) string {
	path, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	return filepath.Join(path, rest)
}

type ReplHistory interface {
	Read(*Repl) error
	Write(*Repl) error
}

type FileReplHist struct {
	Path string
}

func (h *FileReplHist) Read(r *Repl) error {
	if h.Path == "" {
		return nil
	}
	f, err := os.Open(h.Path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = r.ReadHistory(f)
	return err
}
func (h *FileReplHist) Write(r *Repl) error {
	if h.Path == "" {
		return nil
	}
	dir := filepath.Dir(h.Path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	f, err := os.Create(h.Path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = r.WriteHistory(f)
	return err
}

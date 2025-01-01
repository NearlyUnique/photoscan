package main

import (
	"os"
)

type FileOpener struct {
	Filename string
	file     *os.File
}

func (f *FileOpener) Open() (*os.File, error) {
	if f.file == nil {
		return os.Open(f.Filename)
	}
	return f.file, nil
}

func (f *FileOpener) Close() error {
	if f.file != nil {
		return f.file.Close()
	}
	return nil
}

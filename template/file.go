package template

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FilePipe struct {
	Name        string
	Root        string
	d           *bytes.Buffer
	info        os.FileInfo
	files       []string
	destination string
}

func (f *FilePipe) Read() io.Reader {
	return f.d
}

func (f *FilePipe) Write(b []byte) (int, error) {
	return f.d.Write(b)
}
func (f *FilePipe) FlushAndWrite(b []byte) {
	f.d.Reset()
	f.d.Write(b)
}
func (f *FilePipe) String() string {
	return f.d.String()
}

func (f *FilePipe) init() error {
	info, err := os.Stat(f.Name)
	if err != nil {
		return err
	}
	f.info = info
	if info.IsDir() {
		err = f.getAllFiles()
		if err != nil {
			return err
		}
		f.Root = f.Name
		return nil
	}
	b, err := ioutil.ReadFile(f.Name)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	return err
}

func (f *FilePipe) AllFiles() []string {
	return f.files
}
func (f *FilePipe) getAllFiles(args ...string) error {
	var base string
	if len(args) > 0 {
		base = args[0]
	} else {
		base = f.Name
	}
	return filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		f.files = append(f.files, path)
		return nil
	})
}
func newFilePipe(name string) *FilePipe {
	return &FilePipe{
		Name: name,
		d:    &bytes.Buffer{},
	}
}

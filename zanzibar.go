package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gernest/zanzibar/template"
)

func run(t *template.Template, cmds ...string) (err error) {
	if len(cmds) > 0 {
		for _, v := range cmds {
			err = t.ExecuteTemplate(os.Stdout, v, nil)
		}
		return
	}
	return t.Execute(&writer{os.Stdout}, nil)
}

func main() {
	var (
		files    listArg
		commands listArg
		base     = "ZANZIBAR"
	)
	flag.Var(&files, "s", "a comma separated list of files to be used as build scripts")
	flag.Var(&commands, "x", "comma separated  list of commands to be run")
	flag.Parse()
	var f []string
	f = append(f, base)
	f = append(f, files...)
	tmpl, err := template.ParseFiles(f...)
	if err != nil {
		log.Fatal(err)
	}
	err = run(tmpl, commands...)
	if err != nil {
		fmt.Fprintf(os.Stdout, "problem running commands %s", err)
	}
}

type listArg []string

func (l *listArg) String() string {
	return fmt.Sprint(*l)
}

func (l *listArg) Set(value string) error {
	if len(*l) > 0 {
		return errors.New("argument list flag already set")
	}
	for _, s := range strings.Split(value, ",") {
		*l = append(*l, s)
	}
	return nil
}

type writer struct {
	out io.Writer
}

func (w *writer) Write(p []byte) (int, error) {
	s := bytes.TrimSpace(p)
	if len(s) > 0 {
		buf := &bytes.Buffer{}
		buf.WriteString("=> ")
		buf.Write(s)
		buf.WriteString("\n")
		return w.out.Write(buf.Bytes())
	}
	return 0, nil
}

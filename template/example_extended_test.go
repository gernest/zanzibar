package template

import (
	"io/ioutil"
	"log"
	"os"
)

func ExampleFile() {
	b, err := ioutil.ReadFile("testdata/file3.tmpl")
	if err != nil {
		log.Fatalf("executing %s", err)
	}
	tmpl, err := New("test").Parse(string(b))
	if err != nil {
		log.Fatalf("parsing %s", err)
	}
	err = tmpl.Execute(os.Stdout, "doc.go")
	if err != nil {
		log.Fatalf("executing %s", err)
	}

	//Output:
	// testdata/file3.tmpl

}
func ExampleKeyword_file() {
	d := "{{file .}}{{.Name}}{{end}}"
	tmpl, err := New("test").Parse(d)
	if err != nil {
		log.Fatalf("parsing %s", err)
	}
	err = tmpl.Execute(os.Stdout, "testdata")
	if err != nil {
		log.Fatalf("executing %s", err)
	}

}

func ExampleFunc_concat() {
	os.RemoveAll("testdata/concat.tmpl")
	d := `{{file .}}{{concat .  "testdata/concat.tmpl"}}{{end}}`
	tmpl, err := New("test").Parse(d)
	if err != nil {
		log.Fatalf("parsing %s", err)
	}
	err = tmpl.Execute(os.Stdout, "testdata")
	if err != nil {
		log.Fatalf("executing %s", err)
	}
}
func ExampleFunc_copy() {
	os.RemoveAll("testdata/copyout")
	d := `{{file .}}{{copy .  "testdata/copyout"}}{{end}}`
	tmpl, err := New("test").Parse(d)
	if err != nil {
		log.Fatalf("parsing %s", err)
	}
	err = tmpl.Execute(os.Stdout, "testdata/copy")
	if err != nil {
		log.Fatalf("executing %s", err)
	}
}

func ExampleFunc_run() {
	d := `{{run "which" "go"}}`
	tmpl, err := New("test").Parse(d)
	if err != nil {
		log.Fatalf("parsing %s", err)
	}
	err = tmpl.Execute(os.Stdout, "testdata/copy")
	if err != nil {
		log.Fatalf("executing %s", err)
	}
	// Output:
	// /home/gernest/go/bin/go
}

func ExampleFunc_list() {
	d := `{{$dirs:= list "one" "one/two" "one/two/three"}}{{$dirs}}`
	tmpl, err := New("test").Parse(d)
	if err != nil {
		log.Fatalf("parsing %s", err)
	}
	err = tmpl.Execute(os.Stdout, "testdata/copy")
	if err != nil {
		log.Fatalf("executing %s", err)
	}
	// Output:
	// [one one/two one/two/three]

}

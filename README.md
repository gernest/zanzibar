# zanzibar

This is an experimental build tool based on the golang `text/template` library. As a way to learn more about Go programming language.

zanzibar extends the `text/template` libary to include a new `*FileNode`, which makes it easy to deal with file manipulations.

Example the following template code, iterates on the files of the  testdata directory and prints its name.

	{{file "testdata"}}
		{{.Name}}
	{{end}}


Example the following template code saves the content of file `testdata/save.tmpl` to `testdata/halloween.tmpl`

	{{file "testdata/save.tmpl"}}
		{{save . "halloween.tmpl"| .Name}}
	{{end}}

## Installation
	go get github.com/gernest/zanzibar

## Usage
You will need to create a file ZANZIBAR on the root of your project. The file is just a normal golang template.

This is a sample ZANZIBAR file I use on another of my project [aurora](https://github.com/gernest/aurora)

```

{{/* builds aurora */}}

{{/* configurations */}}
    {{$name       := "aurora"}}
    {{$version    := "0.0.1"}}
    {{$public     := "public"}}
    {{$templates  := "templates"}}
    {{$destination:= "builds"}}
    {{$config     := "config"}}
    {{$database   := "db"}}
    {{$cmd        :="cmd/aurora/aurora.go"}}
    {{$buildPath  := printf "%s/%s" $destination $version}}
{{/* end  configuration */}}

{{printf "building %s it might take a while please wait .." $name}}
{{/* setup */}}
    {{/* get all dependencies */}}
    {{run "go" "get" "-t"}}

    {{/* remove any previous builds */}}
    {{clean $destination}}
    {{mkdir $buildPath 0700}}
{{/* end setup */}}

{{/* test */}}
    {{run "go" "test"}}
{{/* end test */}}

{{/* create binary */}}
    {{$bin:=printf "%s/%s" $buildPath $name}}

    {{run "go" "build" "-o" $bin $cmd}}
{{/* end binary*/}}

{{/* assemble */}}
    {{/* prepare database */}}
    {{$dbDir:=printf "%s/%s" $buildPath $database}}
    {{mkdir $dbDir 0700}}

    {{/* copy configurations */}}
    {{$cfg:=printf "%s/%s" $buildPath $config}}
    {{file $config}}
        {{copy . $cfg|ignore}}
    {{end}}

    {{/* copy public files */}}
    {{$pub:=printf "%s/%s" $buildPath $public}}
    {{file $public}}
        {{copy . $pub}}
    {{end}}

    {{/* copy templates */}}
    {{$tmpl:=printf "%s/%s" $buildPath $templates}}
    {{file $templates}}
        {{copy . $tmpl}}
    {{end}}
{{/* end assemble */}}
{{printf "[SUCCESS] built %s version %s" $name $version}}
```

## Running
After you have created the ZANZIBAR file. run zanzibar in the project root

	zanzibar

## Apart from the standard builtin template functions, I added the following.

- concat: concat files.

- copy: copies files or directories.

- run: runs commands.

- mkdir: creates directories.

- list: convert a space separated strings to a slice.

- dest: sets destination for a file.

- save: saves a file to disc.

- clean: deletes directories or files.

- ignore: silence output of the given pipeline.

## Working with files.

The cool part about this effort is the `file` keyword. This behaves just like `with` but instead the context is a file like object of type `FilePipe`.

When the argument to file is a file, it will only pass the `FilePipe` of that file to the children of the node.

	{{file "foo.tmpl"}}
		{{/* here the dot context will refer to FilePipe of file foo.tmpl}}
	{{end}}

When the argument is a  directory, It will iterate over the child nodes for all files in the directory each time passing the current FilePipe.

This is how the file works.

```

func (s *state) walkFile(dot reflect.Value, node *parse.FileNode) {
	defer s.pop(s.mark())
	s.at(node)
	val, _ := indirect(s.evalPipeline(dot, node.Pipe))
	mark := s.mark()
	oneIteration := func(elem reflect.Value) {
		// Set top var (lexically the second if there are two) to the element.
		if len(node.Pipe.Decl) > 0 {
			s.setVar(1, elem)
		}
		s.walk(elem, node.List)
		s.pop(mark)
	}
	switch val.Kind() {
	case reflect.String:
		fileName := val.String()
		fp := newFilePipe(fileName)
		err := fp.init()
		if err != nil {
			log.Fatalf("reading file %v", err)
		}
		if fp.info.IsDir() {
			for _, v := range fp.files {
				nfp := newFilePipe(v)
				err := nfp.init()
				err = nfp.init()
				if err != nil {
					log.Fatalf("reading file %v", err)
				}
				nfp.Root = fp.Root
				oneIteration(reflect.ValueOf(nfp))
			}

		} else {
			s.walk(reflect.ValueOf(fp), node.List)
		}

	}

}
```


## What?
- The `text/template` is huge, I had no guts to also read the whole test suite, so sorry no tests yet(I made changes to the parser test though, to test the `*FileNode`).

- The modified `test/template` is included in this project.

- The golang standard library is like a college to people like me. Its wonderful

## Contributions
contributions are welcome

## License
This project is under the MIT License. See the [LICENSE](LICENCE) file for the full license text.

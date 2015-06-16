package template

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func concat(fp *FilePipe, saveAs string) (*FilePipe, error) {
	f, err := os.OpenFile(saveAs, os.O_WRONLY|os.O_CREATE|os.O_APPEND, fp.info.Mode())
	if err != nil {
		return fp, err
	}
	n, err := f.Write(fp.d.Bytes())
	if err == nil && n < fp.d.Len() {
		return fp, io.ErrShortWrite
	}
	return fp, f.Close()
}

func copyTo(fp *FilePipe, dest string) (string, error) {
	saveTo := getSavingPath(fp, dest)
	info, err := os.Stat(filepath.Dir(fp.Name))
	if err != nil {
		return fmt.Sprintf("copy %s to %s", fp.Name, saveTo), err
	}
	os.MkdirAll(filepath.Dir(saveTo), info.Mode())
	return fmt.Sprintf("copy %s to %s", fp.Name, saveTo), copyFile(fp.Name, saveTo)
}

func copyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourcefile.Close()
	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}
	}
	return
}

func copyDir(src, dest string) (err error) {
	sourceinfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	// create dest dir
	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(src)
	objects, err := directory.Readdir(-1)
	for _, obj := range objects {
		sourcefilepointer := src + "/" + obj.Name()
		destinationfilepointer := dest + "/" + obj.Name()
		if obj.IsDir() {
			// create sub-directories - recursively
			err = copyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = copyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return
}

func getSavingPath(fp *FilePipe, dest string) string {
	if !filepath.IsAbs(dest) {
		return filepath.Join(dest, filepath.Base(fp.Name))
	}
	name := strings.TrimPrefix(fp.Name, fp.Root)
	return filepath.Join(dest, name)
}

func mkdir(path string, mode os.FileMode) (string, error) {
	return fmt.Sprintf("creating %s", path), os.MkdirAll(path, mode)
}
func runCommad(cmd string, args ...string) (string, error) {
	out, err := exec.Command(cmd, args...).Output()
	return string(out), err
}

func sliceEm(args ...string) []string {
	return args
}

func mapEm(args ...string) map[string]string {
	if len(args)%2 != 0 {
		return nil
	}
	m := make(map[string]string)
	for i := 0; i < len(args)+1; i++ {
		m[args[i]] = args[1+1]
	}
	return m
}

func setDestination(fp *FilePipe, dest string) *FilePipe {
	to := getSavingPath(fp, dest)
	fp.destination = to
	return fp
}

func save(fp *FilePipe, args ...string) (*FilePipe, error) {
	var fileName string
	fileName = fp.Name
	if len(args) > 0 {
		fileName = filepath.Join(filepath.Dir(fp.Name), args[0])
	}
	if fp.destination != "" {
		fileName = filepath.Join(fp.destination, filepath.Base(fp.Name))
	}
	info, err := os.Stat(filepath.Dir(fp.Name))
	if err != nil {
		return fp, err
	}
	os.MkdirAll(filepath.Dir(fileName), info.Mode())
	return fp, ioutil.WriteFile(fileName, fp.d.Bytes(), fp.info.Mode())
}

func cleanPath(path string) (string, error) {
	return path, os.RemoveAll(path)
}

func ignoreResult(s string) string {
	return ""
}

package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func WriteDataToFile(path string, data []byte) {
	// make dir if not exists
	if err := os.MkdirAll(filepath.Dir(path), os.ModeDir|0755); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		panic(err)
	}

	return
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		msg := fmt.Sprintf("[utils.CopyFile] non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
		panic(msg)
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			msg := fmt.Sprintf("[utils.CopyFile] non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
			panic(msg)
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}

	copyFileContents(src, dst)
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) {
	in, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		panic(err)
	}

	if err = out.Sync(); err != nil {
		panic(err)
	}

	return
}

// TODO: Consider security: privilege escalation, write to path beyond user rights (what if executable has root?)
func CopyFolderRecursive(srcPath, dst string) {
	cleaned, err := filepath.Abs(filepath.Clean(srcPath))
	if err != nil {
		panic("[utils.CopyFolderRecursive] could not get clean path name")
	}
	baseLen := len(strings.Split(cleaned, "/")) - 1

	var walkFn filepath.WalkFunc

	walkFn = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dir := strings.Join(strings.Split(path, "/")[baseLen:], "/")
		curDst := filepath.Join(dst, dir)

		if info.Mode().IsRegular() {
			copyFileContents(path, curDst)
		} else if info.IsDir() {
			err = os.MkdirAll(curDst, os.ModeDir|0755)
			if err != nil {
				return err
			}
		}

		return nil
	}

	if err := filepath.Walk(srcPath, walkFn); err != nil {
		panic(err)
	}
}

package utils

import (
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "regexp"
    "strings"
)

func TitleFromFilename(filename string) string {
    // get 'title' part
    extension := filepath.Ext(filename)
    title := filename[0 : len(filename)-len(extension)]

    // clean up title by removing possible whitespace or custom wordseparators
    title = strings.Replace(title, "-", " ", -1)
    title = strings.Replace(title, "_", " ", -1)
    title = strings.Trim(title, " ")

    return strings.Title(title)
}

// Post filename must be structured to include a filename prefix of "YYYY-MM-DD"
// What follows after the prefix is optional but it is encouraged to use the title for the post as name.
// Filename will be use as output to site so choose characters wisely.
// Returns error if date part of filname not valid.
func SplitDateAndFilname(filename string) (date, name string) {
    bytes := []byte(filename)
    matcher := regexp.MustCompile("^[0-9]{4}-[0-9]{2}-[0-9]{2}.*")
    if !matcher.Match(bytes) {
        return
    }

    // extract date, title and replace extension
    date = string(bytes[:10])
    name = string(bytes[10:])

    return
}

func WriteFile(path string, data []byte) error {
    // make dir if not exists
    if err := os.MkdirAll(filepath.Dir(path), os.ModeDir|0755); err != nil {
        panic(err)
    }

    return ioutil.WriteFile(path, data, 0644)
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
    sfi, err := os.Stat(src)
    if err != nil {
        return
    }
    if !sfi.Mode().IsRegular() {
        // cannot copy non-regular files (e.g., directories,
        // symlinks, devices, etc.)
        return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
    }
    dfi, err := os.Stat(dst)
    if err != nil {
        if !os.IsNotExist(err) {
            return
        }
    } else {
        if !(dfi.Mode().IsRegular()) {
            return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
        }
        if os.SameFile(sfi, dfi) {
            return
        }
    }

    err = copyFileContents(src, dst)
    return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
    in, err := os.Open(src)
    if err != nil {
        return
    }
    defer in.Close()
    out, err := os.Create(dst)
    if err != nil {
        return
    }
    defer func() {
        cerr := out.Close()
        if err == nil {
            err = cerr
        }
    }()
    if _, err = io.Copy(out, in); err != nil {
        return
    }
    err = out.Sync()
    return
}

// TODO: Consider security: privilege escalation, write to path beyond user rights (what if executable has root?)
func CopyFolderVerbatim(srcPath, dst string) error {

    cleaned, err := filepath.Abs(filepath.Clean(srcPath))
    if err != nil {
        return err
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

            // if !os.SameFile(path, curDst) { // TODO

            err = copyFileContents(path, curDst)
            if err != nil {
                return err
            }
        } else if info.IsDir() {
            err = os.MkdirAll(curDst, os.ModeDir|0755)
            if err != nil {
                return err
            }
        }

        return nil
    }

    return filepath.Walk(srcPath, walkFn)
}

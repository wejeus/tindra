package context

import (
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "regexp"
    "strings"
)

func list(m map[string][]byte) {
    for name, _ := range m {
        fmt.Printf("%s ", name)
    }
}

func show(m map[string][]byte) {
    for name, data := range m {
        fmt.Printf("%s\n----------------------------\n%s\n", name, string(data))
    }
}

// TODO: Consider security: privilege escalation, write to path beyond user rights (what if executable has root?)
func copyFolderVerbatim(srcPath, dst string) error {

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

            fmt.Printf("copy: %s -> %s\n", path, dir)

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

    fmt.Println("Recursive folder copy: " + srcPath)
    return filepath.Walk(srcPath, walkFn)
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

// separates files and folders, skips special entries and hidden entries
func getFilesAndFolders(path string) (specialFiles, verbatimDirs map[string]bool, err error) {
    special := map[string]bool{
        MAIN_CONFIG_FILENAME: true,
        INCLUDES_DIR_NAME:    true,
        LAYOUTS_DIR_NAME:     true,
        POSTS_DIR_NAME:       true,
        DATA_DIR_NAME:        true,
        PLUGINS_DIR_NAME:     true,
        BUILD_DIR_NAME:       true,
    }

    specialFiles = make(map[string]bool)
    verbatimDirs = make(map[string]bool)

    entries, err := ioutil.ReadDir(path)
    if err != nil {
        return
    }

    for _, entry := range entries {
        if special[entry.Name()] || string([]byte(entry.Name())[0:1]) == "." {
            continue
        }

        if entry.IsDir() {
            verbatimDirs[entry.Name()] = true
        } else {
            specialFiles[entry.Name()] = true
        }
    }

    return
}

func applyLayout(layout string, body []byte, layouts map[string]*Layout) ([]byte, error) {
    if len(layout) == 0 {
        return body, nil
    }

    parentLayout := layouts[layout]
    if parentLayout == nil {
        return nil, errors.New("\"" + string(layout) + "\" layout does not exist")
    }

    contentRegexp := regexp.MustCompile("{% content %}")
    return contentRegexp.ReplaceAllLiteral(parentLayout.builtBody, body), nil
}

package context

import (
    "errors"
    "fmt"
    "io/ioutil"
    "path/filepath"
    "regexp"
)

// ALGORITHM
//
// BUILD_INCLUDE(path):
//     if path is read and marked as not generated: abort(cyclic dependency)
//     read file
//     parse file and determine dependencies
//     check each dependency
//         if read: continue
//         else: BUILD_INCLUDE(dependency path)

type Includes map[string]Include

type Include struct {
    Body      []byte
    Generated bool
}

/* path is absolute dir to folder of includes */
func (i Includes) ReadDir(path string) error {
    fmt.Printf("reading includes from: %s\n", path)

    defer func() {
        if r := recover(); r != nil {
            fmt.Println("<PANIC!> ", r)
        }
    }()

    // TODO: Find better way to GLOB .html files

    fileInfos, err := ioutil.ReadDir(path)
    if err != nil {
        return err
    }

    for _, fileInfo := range fileInfos {
        ext := filepath.Ext(fileInfo.Name())

        if len(ext) == 0 || ext[1:] != "html" {
            continue
        }

        if err := i.buildInclude(path, fileInfo.Name()); err != nil {
            return err
        }
    }

    return err
}

func (i Includes) buildInclude(path string, filename string) error {

    uri := filepath.Join(path, filename)

    fmt.Println("parsing: " + filename)

    if i.has(filename) && i.get(filename).Generated == false {
        // fmt.Println("cyclic dependency detected in: " + filename)
        return errors.New("cyclic dependency detected in: " + filename)
    }

    fileContent, err := ioutil.ReadFile(uri)
    if err != nil {
        return err
    }

    i[filename] = Include{Body: nil, Generated: false}

    pattern := regexp.MustCompile("{% include .* %}")

    for pattern.Match(fileContent) {
        fileContent = pattern.ReplaceAllFunc(fileContent, func(match []byte) []byte {

            nameRegexp := regexp.MustCompile(`[a-zA-Z0-9]*\.[a-zA-Z0-9]*`)
            includeName := nameRegexp.Find(match)

            if err = i.buildInclude(path, string(includeName)); err != nil {
                panic(err)
            }

            return i.get(string(includeName)).Body
        })
    }

    i[filename] = Include{Body: fileContent, Generated: true}

    return nil
}

func (i Includes) has(s string) bool {
    _, ok := i[s]
    return ok
}

func (i Includes) get(s string) Include {
    return i[s]
}

func (i Includes) printAllIncludes() {
    for name, include := range i {
        fmt.Printf("\n---------------\n%s\n---------------\n%s\n", name, string(include.Body))
    }
}

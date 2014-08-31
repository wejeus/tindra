package models

import (
    "bytes"
    "errors"
    "fmt"
    "io/ioutil"
    "path/filepath"
    "regexp"
)

// Includes is parsed in a tree bottom up order. If a node rely on a another node not
// previously seen above it this node is parsed recurivly and the used.
// if a node encounters a node above it which have been seen previously but not marked
// as completly generated it means a dependency cycle have been detected.

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
func (i Includes) ReadDir(path string) (err error) {
    fmt.Printf("reading includes from: %s\n", path)

    // TODO: Find better way to GLOB .html files

    fileInfos, err := ioutil.ReadDir(path)
    if err != nil {
        return err
    }

    defer func() {
        if r := recover(); r != nil {
            switch x := r.(type) {
            case string:
                err = errors.New(x)
            case error:
                err = x
            default:
                err = errors.New("unknown panic")
            }
        }
    }()

    for _, fileInfo := range fileInfos {
        filename := fileInfo.Name()
        ext := filepath.Ext(filename)

        if len(ext) == 0 || ext[1:] != "html" {
            continue
        }

        if i.has(filename) {
            continue
        }

        if err := i.buildInclude(path, filename); err != nil {
            return err
        }
    }

    return err
}

func (i Includes) buildInclude(path string, filename string) error {

    if i.has(filename) {
        if i.get(filename).Generated {
            return nil
        } else {
            return errors.New("cyclic dependency detected in: " + filename)
        }
    }

    uri := filepath.Join(path, filename)
    fmt.Println("parsing include: " + filename)

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

    return err
}

func (i Includes) has(s string) bool {
    _, ok := i[s]
    return ok
}

func (i Includes) get(s string) Include {
    return i[s]
}

func (i Includes) String() string {
    var buffer bytes.Buffer
    for name, include := range i {
        s := fmt.Sprintf("---------------\n%s\n---------------\n%s", name, string(include.Body))
        buffer.WriteString(s)
    }

    return buffer.String()
}

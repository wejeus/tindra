package main

import (
    "flag"
    "fmt"
    // "github.com/wejeus/tindra/config"
    "github.com/wejeus/tindra/context"
    // "github.com/wejeus/tindra/models"
    "log"
    "os"
    "path/filepath"
    "syscall"
    // "github.com/russross/blackfriday"
)

// var Debug bool = false // TODO: Add custom logging class

var target string

func init() {
    flag.StringVar(&target, "target", "", "location of site to generate")
    flag.Parse()
}

func main() {

    if len(target) == 0 {
        flag.PrintDefaults()
        os.Exit(int(syscall.EINVAL))
    }

    // TODO: Test if target is valid (make approximation)

    source, err := filepath.Abs(target)
    fmt.Printf("generating: %s\n", target)
    fmt.Printf("path: %s\n\n", source)

    if err != nil {
        log.Fatal("could not get current working directory!")
    }

    site := context.NewContext(source)

    site.Install(filepath.Join(source, "/build"))

    // if err := renderedSite.Install(); err != nil {
    //     fmt.Print(err)
    // }
}

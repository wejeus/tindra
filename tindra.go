package main

import (
    "flag"
    "fmt"
    "github.com/wejeus/tindra/context"
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

    basePath, err := filepath.Abs(target)
    fmt.Printf("BasePath: %s\n", basePath)

    if err != nil {
        log.Fatal("could not get current working directory!")
    }

    _, err = context.NewSite(basePath)
    if err != nil {
        fmt.Print(err)
    }

    // var renderedSite *context.RenderedSite

    // if renderedSite, err = site.Build(); err != nil {
    //     fmt.Print(err)
    // }

    // if err := renderedSite.Install(); err != nil {
    //     fmt.Print(err)
    // }
}

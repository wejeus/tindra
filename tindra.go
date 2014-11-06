package main

import (
	"./context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

// var Debug bool = false // TODO: Add custom logging class

// TODO: Add flag for usage of local relative paths in generation

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

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("error: ", r)
		}
	}()

	site := context.NewContext(source)
	site.Install(filepath.Join(source, "/build"))

	fmt.Printf("\ndone\n")
}

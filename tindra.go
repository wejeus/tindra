package main

import (
	"github.com/wejeus/tindra/context"
	"log"
	// "github.com/russross/blackfriday"
)

func main() {

	site, err := context.NewSite()
	if err != nil {
		log.Print(err)
	}

	err = site.BuildAndInstall()
	if err != nil {
		log.Print(err)
	}
}

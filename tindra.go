package main

import (
	"github.com/wejeus/tindra/context"
	"log"
	// "github.com/russross/blackfriday"
)

func main() {

	_, err := context.NewSite()
	if err != nil {
		log.Fatal("could not create new site!")
	}

}

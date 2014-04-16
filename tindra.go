package main

import (
	"fmt"
	"github.com/wejeus/tindra/context"
	// "github.com/russross/blackfriday"
)

func main() {

	site, err := context.NewSite()
	if err != nil {
		fmt.Print(err)
	}

	err = site.BuildAndInstall()
	if err != nil {
		fmt.Print(err)
	}
}

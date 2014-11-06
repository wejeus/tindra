package models

import (
	"bytes"
	"os"
	"path/filepath"
)

type Pages map[string]*Page

// Reads a directory and all its subdirectories looking for .md and .html files
// and creates a new set of Pages
func NewPages(includes Includes, layouts Layouts, basepath string) Pages {
	pages := make(map[string]*Page, 10)

	var walkFn filepath.WalkFunc

	walkFn = func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		page, err := NewPage(includes, layouts, basepath, path)
		if err != nil {
			panic(err)
		}

		if page != nil {
			pages[page.Permalink()] = page
		}

		return nil
	}

	err := filepath.Walk(basepath, walkFn)
	if err != nil {
		panic(err)
	}

	return pages
}

func (self Pages) String() string {
	var buffer bytes.Buffer
	for _, page := range self {
		buffer.WriteString(page.String())
	}
	return buffer.String()
}

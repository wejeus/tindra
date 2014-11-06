package models

import (
	"bytes"
	"fmt"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"path/filepath"
)

// type Pager interface {
//     Title()
//     Permalink()
//     RelativePath()
//     AbsolutePath()
//     FrontMatter()
// }

type Page struct {
	// Filename string
	// RelDir   string
	// defined as relative path to basepath
	// Example: about/index.md
	Filename    string
	RelDir      string
	FrontMatter FrontMatter
	Template    []byte // template with html, layout/includes injected (if page is markdown)
}

func NewPage(includes Includes, layouts Layouts, basepath, path string) (*Page, error) {
	ext := filepath.Ext(path)
	if ext != ".md" && ext != ".html" {
		return nil, nil
	}

	rel, err := filepath.Rel(basepath, path)
	if err != nil {
		return nil, err
	}

	fmt.Printf("reading page: %s\n", rel)

	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	frontMatter, content, err := extractFrontMatterAndContent(fileContent)
	if err != nil {
		return nil, err
	}

	switch ext {
	case ".md":
		content = blackfriday.MarkdownCommon(content)
	case ".html":
		content, err = includes.ApplyIncludes(content)
	default:
		// ignore others
	}

	if frontMatter.HasLayout() {
		l := layouts.get(frontMatter.Layout)
		content = l.ApplyTo(content)
	}

	page := Page{
		Filename:    filepath.Base(rel),
		RelDir:      filepath.Dir(rel),
		FrontMatter: frontMatter,
		Template:    content,
	}

	return &page, nil
}

func (self Page) Permalink() string {
	htmlFilename := NameWithExtension(self.Filename, ".html")
	return filepath.Join(self.RelDir, htmlFilename)
}

func (self Page) String() string {
	var buffer bytes.Buffer

	s1 := fmt.Sprintf("- Page ---------------------\n"+
		"Filename: %s\n"+
		"RelDir: %s\n"+
		"Permalink(): %s\n"+
		"%s",
		self.Filename, self.RelDir, self.Permalink(), self.FrontMatter)
	s2 := fmt.Sprintf("= Template =================\n%s\n============================\n", self.Template)

	buffer.WriteString(s1)
	buffer.WriteString(s2)

	return buffer.String()
}

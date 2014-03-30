package context

import (
	"errors"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Post struct {
	Title   string
	Parent  *Layout
	Content []byte // Might contain markdown. TODO: Consider possibility to include template code (as addition to markdown)
	Path    string
}

func (site *Site) NewPost(filename string, data []byte) (post *Post) {
	pathFromFilename, titleFromFilename, err := splitFilname(filename)
	if err != nil {
		log.Fatal(err)
	}

	excerpt, body, err := extractExcerptAndBody(data, "\n\n") // TODO: get from site config
	if err != nil {
		log.Fatal(err)
	}

	if excerpt.Layout == "" {
		log.Fatal("error: \"" + filename + "\" needs to declare a layout")
	}
	if site.Layouts[excerpt.Layout] == nil {
		log.Fatal("error: \"" + filename + "\" declared layout does not exist")
	}

	parentLayout := site.Layouts[excerpt.Layout]

	if excerpt.Title == "" {
		excerpt.Title = titleFromFilename
	}

	return &Post{Title: excerpt.Title, Parent: parentLayout, Content: body, Path: pathFromFilename} // TODO: Does this cause a copy upon return?
}

func applyLayout(layout *Layout, content []byte) []byte {
	contentRegexp := regexp.MustCompile("{% content %}")
	return contentRegexp.ReplaceAllLiteral(layout.Data, content)
}

// basepath?
func (p *Post) BuildAndInstall(path string) (err error) {
	html := blackfriday.MarkdownCommon(p.Content)

	rendered := applyLayout(p.Parent, html)

	// render template

	// copy template to build dir under subdir Post.SubPath
	filenameTitle := strings.Replace(p.Title, " ", "_", -1)
	outPath := filepath.Join(path, p.Path)
	filename := filepath.Join(outPath, filenameTitle+".html") // Fixme: better handling

	// Write to disk
	err = os.MkdirAll(outPath, os.ModeDir|0755)
	if err != nil {
		return
	}
	log.Printf("installing: %s\n", filename)
	return ioutil.WriteFile(filename, rendered, 0644)
}

// Post filename must be structured to include a filename prefix of "YYYY-MM-DD"
// What follows after the prefix is optional but it is encouraged to use the title for the post as name.
// If the post content does not contain any excerpt header the filename will be used as title
//
// Returns error if filname not valid.
func splitFilname(filename string) (subpath string, title string, err error) {
	bytes := []byte(filename)
	matcher := regexp.MustCompile("^[0-9]{4}-[0-9]{2}-[0-9]{2}.*")
	if !matcher.Match(bytes) {
		err = errors.New("post filename prefix does not match date convention")
		return
	}

	// extract date, title and extension
	date := string(bytes[:10])
	title = string(bytes[10:])
	extension := filepath.Ext(title)
	title = title[0 : len(title)-len(extension)]

	// clean up title by removing possible whitespace or custom wordseparators
	title = strings.Trim(title, " -")
	title = strings.Title(title)

	subpath = strings.Replace(date, "-", "/", 2)

	return
}

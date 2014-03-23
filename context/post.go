package context

import (
	"errors"
	"github.com/russross/blackfriday"
	"gopkg.in/v1/yaml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// TODO: Should we allow includes in posts?

type Excerpt struct {
	Layout string // might be empty
	Title  string
}

type Post struct {
	Excerpt Excerpt
	Body    string // Might contain markdown. TODO: Consider possibility to include template code (as addition to markdown)
	SubPath string
	RawData []byte
}

func NewPost(filename string, data []byte) (post *Post) {
	log.Printf("Parsing post: %s\n", filename)

	subpath, fileTitle, err := splitFilname(filename)
	if err != nil {
		log.Fatal(err)
	}

	excerpt, body, err := parsePostContent(data, "\n\n") // TODO: get from site config
	if err != nil {
		log.Fatal(err)
	}

	if excerpt.Title == "" {
		excerpt.Title = fileTitle
	}

	return &Post{excerpt, body, subpath, data}
}

func (p *Post) Generate(path string) (err error) {
	// TODO: Maybe include template code in post and render using:
	// tmpl, err := template.ParseFiles("index.template")

	// generate markdown
	rendered := blackfriday.MarkdownCommon([]byte(p.Body))

	// put generated output in template

	// copy template to build dir under subdir Post.SubPath
	title := strings.Replace(p.Excerpt.Title, " ", "_", -1)
	outPath := filepath.Join(path, p.SubPath)
	filename := filepath.Join(outPath, title+".html") // Fixme: better handling
	log.Printf("generating: %s\n", filename)

	// Write to disk
	log.Print(outPath)
	err = os.MkdirAll(outPath, os.ModeDir|0755)
	if err != nil {
		return
	}

	return ioutil.WriteFile(filename, rendered, 0644)
}

// Post filename must be structured to include a filename prefix of "YYYY-MM-DD"
// What follows after the prefix is optional but it is encouraged to use the title for the post as name.
// If the post content does not contain any excerpt header the filename will be used as title
//
// Returns error if filname not valid.
func splitFilname(filename string) (subpath string, title string, err error) {
	bytes := []byte(filename)
	matched, err := regexp.Match("^[0-9]{4}-[0-9]{2}-[0-9]{2}.*", bytes)
	if !matched {
		err = errors.New("post filename prefix does not match date convention")
		return
	}

	date := string(bytes[:10])
	title = strings.Trim(string(bytes[10:]), " -")
	title = strings.Title(title)
	extension := filepath.Ext(title)
	title = title[0 : len(title)-len(extension)]

	subpath = strings.Replace(date, "-", "/", 2)

	return
}

//  By default excerpt is your first paragraph of a post: everything before
//  the first two new lines. Testing if an excerpt is present is simply done by testing if first chars == "---"
//
//      ---
//      title: Example
//      ---
//
//      Second paragraph (post content)
//
func parsePostContent(data []byte, separator string) (excerpt Excerpt, body string, err error) {
	hasExcerpt, err := regexp.Match("^---\n", data) // TODO: add better excerpt regexp: "^---\n.*\n---\n"
	if hasExcerpt {
		// if has excerpt post must consist of 2 parts: (excerpt, body)
		post := strings.SplitN(string(data), separator, 2)
		if len(post) != 2 {
			err = errors.New("could not extract (excerpt, body)")
			return
		}

		err = yaml.Unmarshal([]byte(post[0]), &excerpt)
		body = post[1]
	} else {
		body = string(data)
	}

	return
}

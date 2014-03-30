package context

import (
	"bytes"
	"errors"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// TODO: Make internal structure private and export this to template execution?
type Post struct {
	parent   *Layout
	content  []byte // Contains Markdown
	filename string
	path     string

	// Public to site generation
	Title     string
	Date      string // TODO
	Permalink string
	// TODO: add get excerpt function here
}

func (site *Site) NewPost(filename string, data []byte) (post *Post) {
	date, title, err := splitPostFilname(filename)
	if err != nil {
		log.Fatal(err)
	}

	frontMatter, body, err := extractFrontMatterAndBody(data, "\n\n") // TODO: get from site config
	if err != nil {
		log.Fatal(err)
	}

	if frontMatter.Layout != "" && site.Layouts[frontMatter.Layout] == nil {
		log.Fatal("error: \"" + filename + "\" declared layout that do not exist!")
	}

	parentLayout := site.Layouts[frontMatter.Layout]

	postPath := strings.Replace(date, "-", "/", 2)

	var filenamePath string
	if frontMatter.Title == "" {
		filenamePath = filepath.Join(site.config.BuildPath, title)
	} else {
		filenamePath = filenameFromTitle(postPath, frontMatter.Title)
	}

	return &Post{
		parent:    parentLayout,
		content:   body,
		filename:  filenamePath,
		path:      postPath,
		Title:     frontMatter.Title,
		Date:      date,
		Permalink: filenamePath,
	}
}

// TODO: Change to something fancier (extract first paragraph?)
func (p *Post) Excerpt() string {
	head := p.content[0 : len(p.content)/5]
	return string(blackfriday.MarkdownCommon(head))
}

func applyLayout(layout *Layout, content []byte) []byte {
	if layout == nil {
		return content
	}
	contentRegexp := regexp.MustCompile("{% content %}")
	return contentRegexp.ReplaceAllLiteral(layout.Data, content)
}

func filenameFromTitle(path, title string) string {
	// copy template to build dir under subdir Post.SubPath
	filename := strings.Replace(title, " ", "_", -1)
	return filepath.Join(path, filename+".html") // Fixme: better handling
}

type Page struct {
	Title string // site title
	Posts map[string]*Post
	Post  *Post
}

// Assumes filename and path is absolute
func (p *Post) buildAndInstall(site *Site) (err error) {
	outPath := filepath.Join(site.config.BuildPath, p.path)
	outFile := filepath.Join(site.config.BuildPath, p.filename)
	log.Printf("building: %s\n", outFile)

	// TODO: Only parse markdown for posts
	html := blackfriday.MarkdownCommon(p.content)
	rendered := applyLayout(p.parent, html)
	page := Page{
		Title: p.Title,
		Posts: site.Posts,
		Post:  p,
	}

	// execute template
	t := template.Must(template.New("layout").Parse(string(rendered)))
	var parsed bytes.Buffer
	err = t.Execute(&parsed, page)
	if err != nil {
		return
	}

	// Write to disk
	err = os.MkdirAll(outPath, os.ModeDir|0755)
	log.Print("XXX: " + outPath)
	if err != nil {
		return
	}
	log.Printf("installing: %s\n", outFile)
	return ioutil.WriteFile(outFile, parsed.Bytes(), 0644)
}

// Post filename must be structured to include a filename prefix of "YYYY-MM-DD"
// What follows after the prefix is optional but it is encouraged to use the title for the post as name.
// If the post content does not contain any excerpt header the filename will be used as title
//
// Returns error if filname not valid.
func splitPostFilname(filename string) (date string, title string, err error) {
	bytes := []byte(filename)
	matcher := regexp.MustCompile("^[0-9]{4}-[0-9]{2}-[0-9]{2}.*")
	if !matcher.Match(bytes) {
		err = errors.New("post filename prefix does not match date convention")
		return
	}

	// extract date, title and extension
	date = string(bytes[:10])
	title = string(bytes[10:])
	extension := filepath.Ext(title)
	title = title[0 : len(title)-len(extension)]

	// clean up title by removing possible whitespace or custom wordseparators
	title = strings.Trim(title, " -")
	title = strings.Title(title)

	return
}

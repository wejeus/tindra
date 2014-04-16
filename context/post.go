package context

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

type BuildInfo struct {
	Filename         string // generated output filename
	Folder           string // generated subfolder(s) to use
	AbsoluteBuildDir string // absolut build dir (not including 'Folder')
}

func (b BuildInfo) GetAbsolutePath() string {
	return filepath.Join(b.AbsoluteBuildDir, b.Folder, b.Filename)
}

// TODO: Make internal structure private and export this to template execution?
// type Post struct {
// 	parent          *Layout
// 	content         []byte // Contains Markdown // TODO: Maybe use pointer instead?
// 	renderedContent []byte // content rendered with includes/layouts/data/plugins // TODO: Maybe use pointer instead?
// 	filename        string
// 	FrontMatter     FrontMatter // TODO: Maybe private

// 	buildInfo *BuildInfo
// 	// path     string

// 	// Public to site generation
// 	// Title     string
// 	// Date      string // TODO
// 	// Permalink string
// 	// TODO: add get excerpt function here
// }

type Post struct {
	Name        string // TODO: namechange PostTitle (or just Title everywhere?)
	Body        []byte
	Parent      string // nil if none
	FrontMatter FrontMatter
	Rendered    []byte
}

func (p *Post) Date() string {
	return "DDAAAATE"
}

func readPostsDir(path string, allowedFiles map[string]bool) (posts map[string]*Post, err error) {
	posts = make(map[string]*Post)

	data, err := readFiles(path, allowedFiles)
	if err != nil {
		return
	}

	for name, raw := range data {
		p, postErr := NewPost(name, raw)
		if postErr != nil {
			err = postErr
			return
		}
		posts[name] = p
	}

	return
}

func NewPost(name string, data []byte) (post *Post, err error) {
	post = &Post{Name: name}

	frontMatter, body, err := extractFrontMatterAndBody(data, "\n\n") // TODO: get from site config
	if err != nil {
		return
	}

	var parent string = ""
	if len(frontMatter.Layout) != 0 {
		parent = frontMatter.Layout
	}

	post.Body = body
	post.Parent = parent
	post.FrontMatter = frontMatter

	return
}

// TODO: Better name
type TemplatePage struct {
	Dum       Dummy
	PageTitle string
	Post      Post // TODO should be *
	AllPosts  map[string]*Post
	Data      *Data
	Date      string // TODO
}

type Dummy struct {
	Stuff string
}

// Assumes filename and path is absolute
func (p *Post) Build(site *Site) (err error) {
	fmt.Printf("building: %s\n", p.Name)

	html := blackfriday.MarkdownCommon(p.Body)
	rendered, err := applyLayout(p.Parent, html, site.Layouts)
	// fmt.Println(string(rendered))
	if err != nil {
		return err
	}

	page := TemplatePage{
		Dum:       Dummy{Stuff: "adfsadfasfXXXXXXXX"},
		PageTitle: p.buildTitle(),
		Post:      *p,
		Date:      "asdfasdf",
		AllPosts:  site.Posts,
	}

	// execute template
	t := template.Must(template.New("post").Parse(string(rendered)))
	var parsed bytes.Buffer
	err = t.Execute(&parsed, page)
	if err != nil {
		return err
	}
	p.Rendered = parsed.Bytes()

	fmt.Println(parsed.String())
	if err != nil {
		return
	}

	return
}

// TODO: Change to something fancier (extract first paragraph?)
func (p *Post) Excerpt() string {
	head := p.Body[0 : len(p.Body)/5]
	return string(blackfriday.MarkdownCommon(head))
}

// func (p *Post) Permalink() string {
// 	if p.Rendered == nil {
// 		log.Fatal("need to build first!")
// 	}
// 	return filepath.Join(p.buildInfo.Folder, p.buildInfo.Filename)
// }

// Extend to use different types of permalinks. For now uses format /YYYY/MM/DD/filename.html
func (p *Post) Permalink() (filename, folder string) { // TODO: switch order of returns
	date, filename, err := splitPostFilname(p.Name)
	if err != nil {
		log.Fatal(err)
	}
	folder = strings.Replace(date, "-", "/", 2)
	return
}

func (p *Post) buildTitle() string {
	name, _ := p.Permalink()

	if p.FrontMatter.Title != "" {
		return p.FrontMatter.Title
	}

	// get 'title' part
	extension := filepath.Ext(name)
	title := name[0 : len(name)-len(extension)]

	// clean up title by removing possible whitespace or custom wordseparators
	title = strings.Replace(title, "-", " ", -1)
	title = strings.Replace(title, "_", " ", -1)
	title = strings.Trim(title, " ")

	return strings.Title(title)
}

func (p *Post) Install(site *Site) (err error) {
	// build := p.buildInfo
	// if build == nil {
	// 	log.Fatal("need to build before install!")
	// }
	filename, folder := p.Permalink()

	path := filepath.Join("<get root dir>", folder, filename)
	err = os.MkdirAll(path, os.ModeDir|0755)
	if err != nil {
		return
	}

	fmt.Printf("installing: %s\n", path)
	return ioutil.WriteFile(path, p.Rendered, 0644)
}

func applyLayout(layout string, body []byte, layouts map[string]*Layout) ([]byte, error) {
	if len(layout) == 0 {
		return body, nil
	}

	parentLayout := layouts[layout]
	if parentLayout == nil {
		return nil, errors.New("\"" + string(layout) + "\" layout does not exist")
	}

	contentRegexp := regexp.MustCompile("{% content %}")
	return contentRegexp.ReplaceAllLiteral(parentLayout.builtBody, body), nil
}

// Post filename must be structured to include a filename prefix of "YYYY-MM-DD"
// What follows after the prefix is optional but it is encouraged to use the title for the post as name.
// Filename will be use as output to site so choose characters wisely.
// Returns error if date part of filname not valid.
func splitPostFilname(filename string) (date string, name string, err error) {
	bytes := []byte(filename)
	matcher := regexp.MustCompile("^[0-9]{4}-[0-9]{2}-[0-9]{2}.*")
	if !matcher.Match(bytes) {
		err = errors.New("post filename prefix does not match date convention")
		return
	}

	// extract date, title and replace extension
	date = string(bytes[:10])
	name = string(bytes[10:])

	extension := filepath.Ext(name)
	name = name[0 : len(name)-len(extension)]

	name = strings.Trim(name, " _-")
	name = name + ".html" // TODO: Replace with 'addExtension-ish' lib function

	return
}

package context

import (
	// "bytes"
	"errors"
	"gopkg.in/v1/yaml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	// "text/template"
)

// Implement this on all 'page' types. Used when generating final site to get content
// type GetData interface{}
// TODO: Use go routines for templates "Once constructed, a template may be executed safely in parallel."

type FrontMatter struct {
	// If set, this specifies the layout file to use.
	// Use the layout file name without the file extension.
	// Layout files must be placed in the  'layouts' directory.
	Layout string

	// Title is normaly generated using the filename if page is post.
	// If title is set in FrontMatter this title will be used instead.
	Title string

	// If you need your processed blog post URLs to be something other than
	// the default /year/month/day/title.html then you can set this variable
	// and it will be used as the final URL.
	Permalink string // TODO

	// A date here overrides the date from the name of the post.
	// This can be used to ensure correct sorting of posts. Must have format YYYY-MM-DD.
	Date string // TODO

	// Set to false if you don’t want a specific post to show up when the site is generated.
	// Defaults to true.
	Published string // TODO

	// Similar to categories, one or multiple tags can be added to a post.
	// Also like categories, tags can be specified as a YAML list or a space- separated string.
	Tags []string // TODO

	// TODO
	// Instead of placing posts inside of folders, you can specify one or more
	// categories that the post belongs to. When the site is generated the post
	// will act as though it had been set with these categories normally.
	// Categories (plural key) can be specified as a YAML list or a space-separated string.

	// category
	// categories
}

type Layout struct {
	Name         string
	Data         []byte
	Parent       *Layout
	Dependencies []string //pathnames of includes
}

type Site struct {
	config    *Config
	Generator string
	Includes  map[string][]byte
	Layouts   map[string]*Layout
	Pages     map[string]*Post
	Posts     map[string]*Post
}

func NewSite() (s *Site, err error) {
	log.Print("Generating new site...")

	config := NewConfig()
	config.ReadFromConfigFile() // TODO: should be implicit in NewConfig()

	s = &Site{
		Generator: strings.ToUpper(APP_NAME),
		config:    config,
	}

	s.Includes, err = readFiles(config.IncludesPath, map[string]bool{"html": true, "css": true})
	if err != nil {
		log.Fatal(err)
	}

	s.Layouts, err = s.loadAndPreprocessLayouts(config.LayoutsPath, map[string]bool{"html": true})
	if err != nil {
		log.Fatal(err)
	}
	showLayouts(s.Layouts, false)

	s.Posts, err = s.loadAndPreprocessPosts(config.PostsPath, config.MarkdownExt)
	if err != nil {
		log.Fatal(err)
	}

	s.Pages, err = s.loadAndPreprocessPages(config.BasePath, map[string]bool{"html": true})
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (site *Site) BuildAndInstall() (err error) {
	err = os.RemoveAll(site.config.BuildPath)
	if err != nil {
		return
	}

	err = os.MkdirAll(site.config.BuildPath, os.ModeDir|0755)
	if err != nil {
		return
	}

	for _, p := range site.Posts {
		// PageStruct := PageStruct{Site:s Post:thisPost}
		// err := p.BuildAndInstall(s.config.BuildPath, PageStruct)
		err = p.buildAndInstall(site)
		if err != nil {
			return
		}
	}

	for _, p := range site.Pages {
		// PageStruct := PageStruct{Site:s Post:thisPost}
		// err := p.BuildAndInstall(s.config.BuildPath, PageStruct)
		err = p.buildAndInstall(site)
		if err != nil {
			return
		}
	}

	return
}

func (site *Site) loadAndPreprocessLayouts(directory string, postfixes map[string]bool) (m map[string]*Layout, err error) {
	var data map[string][]byte
	m = make(map[string]*Layout)

	data, err = readFiles(directory, postfixes)
	if err != nil {
		return
	}

	for layoutName, raw := range data {
		m[layoutName] = &Layout{Name: layoutName, Data: raw}
	}

	for layoutName, layoutStruct := range m {
		frontMatter, body, extractErr := extractFrontMatterAndBody(layoutStruct.Data, "\n\n")
		if extractErr != nil {
			err = extractErr
			return
		}

		if len(frontMatter.Layout) != 0 {
			if m[frontMatter.Layout] == nil {
				// this layout has en frontMatter but refering to a layout that does not exist
				err = errors.New("can't find layout dependency: " + frontMatter.Layout)
				return
			}

			m[layoutName].Data = body
		}

		m[layoutName].Parent = m[frontMatter.Layout]
	}

	for name, layout := range m {

		// TODO: check for existens of {% content %} tag (only one per layout allowed)

		parentLayout := layout.Parent
		content := layout.Data
		for parentLayout != nil {

			// TODO check dependecy graph for existens and circular dependencies

			contentRegexp := regexp.MustCompile("{% content %}")
			content = contentRegexp.ReplaceAllLiteral(parentLayout.Data, content)
			m[name].Data = content

			// iterate to next
			parentLayout = parentLayout.Parent
		}

		// TODO: Regexp file extention should match only extensions defined in config
		includeRegexp := regexp.MustCompile("{% include .* %}")
		m[name].Data = includeRegexp.ReplaceAllFunc(m[name].Data, func(match []byte) []byte {
			nameRegexp := regexp.MustCompile(`[a-zA-Z0-9]*\.[a-zA-Z0-9]*`)
			includeName := nameRegexp.Find(match)
			if includeName == nil || site.Includes[string(includeName)] == nil {
				err = errors.New("invalid include name or include does not exist")
			}

			layout.Dependencies = append(layout.Dependencies, string(includeName))
			return site.Includes[string(includeName)]
		})

	}
	return
}

func list(m map[string][]byte) {
	for name, _ := range m {
		log.Printf("%s ", name)
	}
}

func show(m map[string][]byte) {
	for name, data := range m {
		log.Printf("%s\n----------------------------\n%s\n", name, string(data))
	}
}

func showLayouts(m map[string]*Layout, showContents bool) {
	for name, layout := range m {
		log.Printf("Layout: %s\n", name)

		var parent string
		if layout.Parent == nil {
			parent = "<none>"
		} else {
			parent = layout.Parent.Name
		}
		log.Printf("Parent: %s\n", parent)

		deps := ""
		for _, includeName := range layout.Dependencies {
			deps = deps + includeName + ", "
		}
		if len(deps) == 0 {
			deps = "<none>"
		}
		log.Printf("Includes: %s\n", deps)

		if showContents {
			log.Printf("Content:\n%s\n", layout.Data)
		}
	}
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
// if file does not contain any frontMatter only the body will be returned and all other values will be nil
func extractFrontMatterAndBody(data []byte, separator string) (frontMatter FrontMatter, body []byte, err error) {
	// TODO: Change regexp to use .MustCompile
	hasFrontMatter, err := regexp.Match("^---\n", data) // TODO: add better frontMatter regexp: "^---\n.*\n---\n"
	if hasFrontMatter {
		// if has frontMatter post must consist of 2 parts: (frontMatter, body)
		post := strings.SplitN(string(data), separator, 2)
		if len(post) != 2 {
			err = errors.New("could not extract (frontMatter, body)")
			return
		}

		// TODO: Check that frontMatter only contains one value for key 'layout'
		err = yaml.Unmarshal([]byte(post[0]), &frontMatter)
		body = []byte(post[1])
	} else {
		body = data
	}

	return
}

// func installDir(directory string) {
// 	os.Wal
// 	// Write to disk
// 	err = os.MkdirAll(outPath, os.ModeDir|0755)
// 	if err != nil {
// 		return
// 	}
// 	log.Printf("installing: %s\n", outFile)
// 	return ioutil.WriteFile(outFile, parsed.Bytes(), 0644)
// }

// TODO: Add possible regex match on filename instead?
// TODO: Also read subfolders and prepend subfolder name to key
// TODO: function name change: readFilesRaw?
// TODO: Implement recursive read of dirs
// Reads all files with with postfix set to true in postfixes map.
// A new map of {parentPaths/filename, content} is returned or error if something goes bananas.
func readFiles(directory string, postfixes map[string]bool) (m map[string][]byte, err error) {
	log.Printf("Reading directory: %s\n", directory)
	m = make(map[string][]byte)

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return
	}

	for _, f := range files {
		extension := filepath.Ext(f.Name())
		if len(extension) != 0 && postfixes[extension[1:]] {
			uri := filepath.Join(directory, f.Name())
			data, err := ioutil.ReadFile(uri)
			if err != nil {
				return m, err
			}
			m[f.Name()] = data
		}
	}

	return
}

func (site *Site) loadAndPreprocessPosts(directory string, postfixes map[string]bool) (posts map[string]*Post, err error) {
	files, err := readFiles(directory, postfixes)
	if err != nil {
		return
	}

	posts = make(map[string]*Post)
	for filename, data := range files {
		post := site.NewPost(filename, data)
		posts[filename] = post
	}

	return
}

func (site *Site) loadAndPreprocessPages(directory string, postfixes map[string]bool) (posts map[string]*Post, err error) {
	files, err := readFiles(directory, postfixes)
	if err != nil {
		return
	}

	posts = make(map[string]*Post)
	for filename, data := range files {
		log.Print("parsing page: " + filename)
		post := site.NewPage(filename, data)
		posts[filename] = post
	}

	return
}

func (site *Site) NewPage(filename string, data []byte) (post *Post) {
	frontMatter, body, err := extractFrontMatterAndBody(data, "\n\n") // TODO: get from site config
	if err != nil {
		log.Fatal(err)
	}

	if frontMatter.Layout != "" && site.Layouts[frontMatter.Layout] == nil {
		log.Fatal("error: \"" + filename + "\" declared layout that do not exist!")
	}

	title := site.config.Name
	if len(frontMatter.Title) != 0 {
		title = frontMatter.Title
	}

	parentLayout := site.Layouts[frontMatter.Layout]

	return &Post{
		parent:   parentLayout,
		content:  body,
		filename: filename,
		path:     ".",

		Title: title,
	} // TODO: Does this cause a copy upon return?
}

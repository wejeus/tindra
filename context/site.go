package context

import (
	// "bytes"
	"errors"
	"gopkg.in/v1/yaml"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	// "text/template"
)

// Implement this on all 'page' types. Used when generating final site to get content
// type GetData interface{}
// TODO: Use go routines for templates "Once constructed, a template may be executed safely in parallel."

type Excerpt struct {
	// If set, this specifies the layout file to use.
	// Use the layout file name without the file extension.
	// Layout files must be placed in the  'layouts' directory.
	Layout string

	// If you need your processed blog post URLs to be something other than
	// the default /year/month/day/title.html then you can set this variable
	// and it will be used as the final URL.
	Permalink string

	// Title is normaly generated using the filename of post.
	// If title is set in excerpt this title will be used instead.
	Title string

	// A date here overrides the date from the name of the post.
	// This can be used to ensure correct sorting of posts. Must have format YYYY-MM-DD.
	Data string

	// Set to false if you don’t want a specific post to show up when the site is generated.
	// Defaults to true.
	Published string

	// Similar to categories, one or multiple tags can be added to a post.
	// Also like categories, tags can be specified as a YAML list or a space- separated string.
	Tags []string

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
	Posts     map[string]*Post
}

func NewSite() (s *Site, err error) {
	config := NewConfig()
	config.ReadFromConfigFile()

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

	for _, p := range s.Posts {
		// PageStruct := PageStruct{Site:s Post:thisPost}
		// err := p.BuildAndInstall(s.config.BuildPath, PageStruct)
		err := p.BuildAndInstall(s.config.BuildPath)
		if err != nil {
			log.Fatal(err)
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
		excerpt, body, extractErr := extractExcerptAndBody(layoutStruct.Data, "\n\n")
		if extractErr != nil {
			err = extractErr
			return
		}

		if len(excerpt.Layout) != 0 {
			if m[excerpt.Layout] == nil {
				// this layout has en excerpt but refering to a layout that does not exist
				err = errors.New("can't find layout dependency: " + excerpt.Layout)
				return
			}

			m[layoutName].Data = body
		}

		m[layoutName].Parent = m[excerpt.Layout]
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

// TODO: This will most likely be a chain when post can use layouts which in turn use layouts...
// type Page struct {
// 	s       *Site
// 	parent  string
// 	Content string
// }

// func (s *Site) renderLayouts() (renderedLayouts map[string][]byte, err error) {

// 	templateDataStruct := make(map[string]Page)

// 	for layoutName, layoutData := range s.RawLayouts {
// 		excerpt, body, extractErr := extractExcerptAndBody(layoutData, "\n\n")
// 		if extractErr != nil {
// 			err = extractErr
// 			return
// 		}

// 		// First render each layout separatly (excluding excerpt) then if has excerpt use defined layout as parent

// 		// t := template.Must(template.New("layout").Parse(string(body)))
// 		// var parsed bytes.Buffer
// 		// p := Page{s: s}
// 		// t.Execute(&parsed, p)

// 		if len(excerpt.Layout) != 0 {
// 			parent := s.RawLayouts[excerpt.Layout]
// 			templateDataStruct[layoutName] = Page{s: s, parent: string(parent), Content: string(body)}

// 			// t := template.Must(template.New("layout").Parse(string(parent)))
// 			// var doc bytes.Buffer

// 			// p := Page{}
// 			// t.Execute(&doc, p)
// 			// renderedLayouts[layout] = doc.Bytes()
// 		} else {
// 			templateDataStruct[layoutName] = Page{s: s, parent: "", Content: string(body)}
// 		}
// 	}

// 	renderedLayouts = make(map[string][]byte)
// 	for layoutName, page := range templateDataStruct {

// 		if len(page.parent) != 0 {
// 			t := template.Must(template.New(layoutName).Parse(page.parent))
// 			var parsed bytes.Buffer
// 			t.Execute(&parsed, page)
// 			renderedLayouts[layoutName] = parsed.Bytes()
// 		}
// 	}

// 	return
// }

//  By default excerpt is your first paragraph of a post: everything before
//  the first two new lines. Testing if an excerpt is present is simply done by testing if first chars == "---"
//
//      ---
//      title: Example
//      ---
//
//      Second paragraph (post content)
//

// if file does not contain any excerpt only the body will be returned and all other values will be nil
func extractExcerptAndBody(data []byte, separator string) (excerpt Excerpt, body []byte, err error) {
	// TODO: Change regexp to use .MustCompile
	hasExcerpt, err := regexp.Match("^---\n", data) // TODO: add better excerpt regexp: "^---\n.*\n---\n"
	if hasExcerpt {
		// if has excerpt post must consist of 2 parts: (excerpt, body)
		post := strings.SplitN(string(data), separator, 2)
		if len(post) != 2 {
			err = errors.New("could not extract (excerpt, body)")
			return
		}

		// TODO: Check that excerpt only contains one value for key 'layout'
		err = yaml.Unmarshal([]byte(post[0]), &excerpt)
		body = []byte(post[1])
	} else {
		body = data
	}

	return
}

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

// TODO: use reflection! http://golang.org/pkg/reflect/
// func (s *Site) Config(key string) (string, error) {
// 	// TODO: Read and return configuration value for key
// 	return "not_set", nil
// }

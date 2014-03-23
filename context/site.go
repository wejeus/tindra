package context

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	// "blackfriday"
)

// Implement this on all 'page' types. Used when generating final site to get content
type GetData interface{}

type Site struct {
	config    *Config
	Generator string
	Includes  map[string]string // TODO: Modify to hold list of type "Include"?
	Layouts   map[string]string // TODO: Modify to hold list of type "Layout"?
	Posts     map[string]*Post
}

func NewSite() (s *Site, err error) {
	config := NewConfig()
	config.ReadFromConfigFile()

	s = &Site{
		Generator: strings.ToUpper(APP_NAME),
		config:    config,
	}
	// s.Includes = make(map[string]string)

	// s.readIncludes()
	// s.readLayouts()
	s.parsePosts()

	for _, p := range s.Posts {
		err := p.Generate(s.config.BuildPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	return
}

func (s *Site) Generate() {

}

func (s *Site) showIncludes() {
	for k, v := range s.Includes {
		log.Printf("%s\n%s\n", k, v)
	}
}

func readFiles(directory string, postfixMatch map[string]bool) (m map[string][]byte, err error) {

	// TODO: Add possible regex match on filename

	log.Printf("Reading directory: %s\n", directory)

	m = make(map[string][]byte)

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return
	}

	for _, f := range files {
		// if filepath.Ext(f.Name()) == ".html" {
		extension := filepath.Ext(f.Name())
		if len(extension) != 0 && postfixMatch[extension[1:]] {
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

// func (s *Site) readIncludes() {
// 	m, err := readFiles(s.config.IncludesPath, ".html")
// 	if err != nil {
// 		log.Fatal("could not read includes!")
// 	} else {
// 		s.Includes = m
// 	}
// }

// func (s *Site) readLayouts() {
// 	m, err := readFiles(s.config.LayoutsPath, ".html")
// 	if err != nil {
// 		log.Fatal("could not read layouts!")
// 	} else {
// 		s.Layouts = m
// 	}
// }

func (s *Site) parsePosts() {
	files, err := readFiles(s.config.PostsPath, s.config.MarkdownExt)
	if err != nil {
		log.Fatal("could not read posts!")
	}

	posts := make(map[string]*Post)
	for filename, data := range files {
		post := NewPost(filename, data)
		posts[filename] = post
	}

	s.Posts = posts
}

// TODO: use reflection! http://golang.org/pkg/reflect/
// func (s *Site) Config(key string) (string, error) {
// 	// TODO: Read and return configuration value for key
// 	return "not_set", nil
// }

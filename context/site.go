package context

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

type Site struct {
	config    *Config
	Generator string
	Includes  map[string]string
	Layouts   map[string]string
}

func NewSite() (s *Site, err error) {
	config := NewConfig()
	config.ReadFromConfigFile()

	s = &Site{
		Generator: strings.ToUpper(APP_NAME),
		config:    config,
	}
	// s.Includes = make(map[string]string)

	s.readIncludes()
	s.readLayouts()
	return
}

func (s *Site) showIncludes() {
	for k, v := range s.Includes {
		log.Printf("%s\n%s\n", k, v)
	}
}

func readFiles(path string, postfix string) (m map[string]string, err error) {
	log.Printf("Reading \"%s\" from: %s\n", postfix, path)

	m = make(map[string]string)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".html" {
			uri := filepath.Join(path, f.Name())
			data, err := ioutil.ReadFile(uri)
			if err != nil {
				return m, err
			}
			m[f.Name()] = string(data)
		}
	}

	return
}

func (s *Site) readIncludes() {
	m, err := readFiles(s.config.IncludesPath, ".html")
	if err != nil {
		log.Fatal("could not read includes!")
	} else {
		s.Includes = m
	}
}

func (s *Site) readLayouts() {
	m, err := readFiles(s.config.LayoutsPath, ".html")
	if err != nil {
		log.Fatal("could not read layouts!")
	} else {
		s.Layouts = m
	}
}

func (s *Site) parsePosts() {

}

// TODO: use reflection! http://golang.org/pkg/reflect/
// func (s *Site) Config(key string) (string, error) {
// 	// TODO: Read and return configuration value for key
// 	return "not_set", nil
// }

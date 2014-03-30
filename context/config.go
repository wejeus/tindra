package context

import (
	"flag"
	"gopkg.in/v1/yaml"
	"io/ioutil"
	"log"
	"path/filepath"
)

func init() {
	flag.Parse()
	// TODO: check if target is valid Tindra site URL (check for existens of required dirs/files)
}

type Config struct {
	// Pygments         bool
	// Host             string
	// Port             int
	// BaseDir          string
	Name             string
	ExcerptSeparator string
	MarkdownExt      map[string]bool
	BasePath         string
	IncludesPath     string
	LayoutsPath      string
	PostsPath        string
	BuildPath        string
}

func NewConfig() *Config {
	target := ""
	if len(flag.Args()) == 1 {
		target = flag.Args()[0] + "/"
	} else if len(flag.Args()) > 1 {
		log.Fatal("can only generate one target!")
	}

	// TODO: Trim to pretty path (remove unnecessary "./" and similar)

	basePath, err := filepath.Abs(target)
	log.Printf("BasePath: %s\n", basePath)

	if err != nil {
		log.Fatal("could not get current working directory!")
	}

	// TODO: Read path of executable? Maybe needed for something later?
	// processWorkDir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	// TODO: preped all dirs with basepath directly

	defaults := Config{
		// Host:             "0.0.0.0",
		// Port:             4000,
		MarkdownExt: map[string]bool{
			"markdown": true,
			"mkdown":   true,
			"mkdn":     true,
			"mkd":      true,
			"md":       true,
		},
		BasePath:     basePath,
		IncludesPath: filepath.Join(basePath, "includes"),
		LayoutsPath:  filepath.Join(basePath, "layouts"),
		PostsPath:    filepath.Join(basePath, "posts"),
		BuildPath:    filepath.Join(basePath, "_build"),
	}
	return &defaults
}

/**
If config file could not be read a default config and an error is returned.
Its up to callee to decide if we should continue or not.
*/
func (c *Config) ReadFromConfigFile() (err error) {
	uri := filepath.Join(c.BasePath, MAIN_CONFIG_FILENAME)
	log.Printf("Reading config file: %s\n", uri)

	content, err := ioutil.ReadFile(uri)
	if err != nil {
		log.Printf("could not read configuration file: %s", uri)
	}

	yaml.Unmarshal(content, c)

	return err
}

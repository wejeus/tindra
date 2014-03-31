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
}

type Config struct {
	// FrontMatterSeparator string // TODO: Add!
	Name        string
	Debug       bool
	MarkdownExt map[string]bool
	BasePath    string
	// IncludesPath string
	// LayoutsPath  string
	// PostsPath    string
	// DataPath     string
	// PluginsPath  string
	// BuildPath    string
}

func NewConfig() *Config {
	target := ""
	if len(flag.Args()) == 1 {
		target = flag.Args()[0] + "/"
	} else if len(flag.Args()) > 1 {
		log.Fatal("can only generate one target!")
	}

	basePath, err := filepath.Abs(target)
	log.Printf("BasePath: %s\n", basePath)

	if err != nil {
		log.Fatal("could not get current working directory!")
	}

	defaults := Config{
		Name:  "Tindra ver. " + VERSION + " " + TAGLINE,
		Debug: DEBUG,
		MarkdownExt: map[string]bool{
			"markdown": true,
			"mkdown":   true,
			"mkdn":     true,
			"mkd":      true,
			"md":       true,
		},
		BasePath: basePath,
		// IncludesPath: filepath.Join(basePath, INCLUDES_DIR_NAME),
		// LayoutsPath:  filepath.Join(basePath, LAYOUTS_DIR_NAME),
		// PostsPath:    filepath.Join(basePath, POSTS_DIR_NAME),
		// DataPath:     filepath.Join(basePath, DATA_DIR_NAME),
		// PluginsPath:  filepath.Join(basePath, PLUGINS_DIR_NAME),
		// BuildPath:    filepath.Join(basePath, BUILD_DIR_NAME),
	}
	return &defaults
}

// If config file could not be read a default config and an error is returned.
// Its up to callee to decide if we should continue or not.
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

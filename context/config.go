package context

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"path/filepath"
)

// var Debug bool = false // TODO: Add custom logging class

func init() {
	flag.Parse()
}

type Config struct {
	// FrontMatterSeparator string // TODO: Add!
	SiteName    string
	MarkdownExt map[string]bool
	basePath    string // should not be able to read from file.
}

func NewConfig() *Config {
	target := ""
	if len(flag.Args()) == 1 {
		target = flag.Args()[0] + "/"
	} else if len(flag.Args()) > 1 {
		log.Fatal("can only generate one target!")
	}

	basePath, err := filepath.Abs(target)
	fmt.Printf("BasePath: %s\n", basePath)

	if err != nil {
		log.Fatal("could not get current working directory!")
	}

	defaults := Config{
		SiteName: "Tindra ver. " + VERSION + " " + TAGLINE,
		MarkdownExt: map[string]bool{
			"markdown": true,
			"mkdown":   true,
			"mkdn":     true,
			"mkd":      true,
			"md":       true,
		},
		basePath: basePath,
	}
	return &defaults
}

func (c *Config) getAbsBuildPath() string {
	return filepath.Join(c.basePath, BUILD_DIR_NAME)
}

func (c *Config) prependAbsPath(dir string) string {
	return filepath.Join(c.basePath, dir)
}

func (c *Config) prependAbsBuildPath(dir string) string {
	return filepath.Join(c.basePath, BUILD_DIR_NAME, dir)
}

// If config file could not be read a default config and an error is returned.
// Its up to callee to decide if we should continue or not.
func (c *Config) ReadFromConfigFile() (err error) {
	uri := filepath.Join(c.basePath, MAIN_CONFIG_FILENAME)
	fmt.Printf("Reading config file: %s\n", uri)

	content, err := ioutil.ReadFile(uri)
	if err != nil {
		fmt.Printf("could not read configuration file: %s", uri)
	}

	yaml.Unmarshal(content, c)

	return err
}

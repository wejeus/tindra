package config

import (
	"fmt"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"path/filepath"
)

type Config struct {
	SiteName    string
	MarkdownExt map[string]bool
}

func NewDefaultConfig() *Config {
	defaults := Config{
		SiteName: "Tindra ver. " + VERSION + " " + TAGLINE,
		MarkdownExt: map[string]bool{
			"markdown": true,
			"mkdown":   true,
			"mkdn":     true,
			"mkd":      true,
			"md":       true,
		},
	}
	return &defaults
}

// If config file could not be read a default config and an error is returned.
// Its up to callee to decide if we should continue or not.
func (c *Config) ReadFromConfigFile(path string) (err error) {
	return c.ReadFromConfigFileNamed(path, "")
}

func (c *Config) ReadFromConfigFileNamed(path, configFile string) (err error) {
	var uri string
	if len(configFile) == 0 {
		uri = filepath.Join(path, MAIN_CONFIG_FILENAME)
	} else {
		uri = filepath.Join(path, configFile)
	}

	fmt.Printf("Reading config file: %s\n", uri)

	content, err := ioutil.ReadFile(uri)
	if err != nil {
		fmt.Printf("could not read configuration file: %s", uri)
	}

	yaml.Unmarshal(content, c)

	return err
}

package context

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Include struct {
	Name      string
	Body      []byte
	Includes  []string //pathnames of sub includes (TODO: sub includes not implemented)
	builtBody []byte   // TODO: namechange -> 'PreProcessed' (may contain template tags)
}

func readIncludesDir(path string) (includes map[string]*Include, err error) {
	allowedFiles := map[string]bool{"html": true} // TODO: Add css/js here?
	includes = make(map[string]*Include)

	data, err := readFiles(path, allowedFiles)
	if err != nil {
		return
	}

	for name, raw := range data {
		inc := NewInclude(name, raw)
		includes[name] = inc
	}

	return
}

func NewInclude(name string, body []byte) *Include {
	include := Include{Name: name, Body: body}

	includeRegexp := regexp.MustCompile("{% include .* %}")
	_ = includeRegexp.ReplaceAllFunc(body, func(match []byte) []byte {
		nameRegexp := regexp.MustCompile(`[a-zA-Z0-9]*\.[a-zA-Z0-9]*`)
		includeName := nameRegexp.Find(match)
		include.Includes = append(include.Includes, string(includeName))
		return match
	})

	return &include
}

func BuildIncludes(includes map[string]*Include) error { // add 'data', 'plugins' as params
	for name, this := range includes {
		fmt.Println("Building: " + name)

		includeRegexp := regexp.MustCompile("{% include .* %}")
		builtBody := this.Body
		err := []string{}

		for includeRegexp.Match(builtBody) {
			builtBody = includeRegexp.ReplaceAllFunc(builtBody, func(match []byte) []byte {
				nameRegexp := regexp.MustCompile(`[a-zA-Z0-9]*\.[a-zA-Z0-9]*`)
				includeName := nameRegexp.Find(match)
				fmt.Println("match: " + string(includeName))
				if len(includeName) == 0 || includes[string(includeName)] == nil {
					err = append(err, "\""+string(includeName)+"\" invalid include name or include does not exist")
					return match
				}
				return includes[string(includeName)].Body
			})

			// check for error here so we can break the for loop
			if len(err) != 0 {
				return errors.New("could not build inludes\n\t" + strings.Join(err, "\n\t"))
			}
		}

		this.builtBody = builtBody
	}

	return nil
}

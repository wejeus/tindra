package context

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Layout struct {
	Name      string
	Body      []byte
	Parent    string // nil if root layout
	Includes  []string
	builtBody []byte // TODO: namechange -> 'PreProcessed' (may contain template tags)
}

func readLayoutsDir(path string) (layouts map[string]*Layout, err error) {
	allowedFiles := map[string]bool{"html": true}
	layouts = make(map[string]*Layout)

	data, err := readFiles(path, allowedFiles)
	if err != nil {
		return
	}

	for name, raw := range data {
		l, layoutErr := NewLayout(name, raw)
		if layoutErr != nil {
			err = layoutErr
			return
		}
		layouts[name] = l
	}

	return
}

// Assumes includes have already been built
func BuildLayouts(layouts map[string]*Layout, includes map[string]*Include) error { // add 'data', 'plugins' as params
	for name, this := range layouts {
		fmt.Println("Building: " + name)

		parent := layouts[this.Parent]
		if len(this.Parent) != 0 && parent == nil {
			return errors.New("in file: \"" + name + "\" layout \"" + this.Parent + "\" does not exist")
		}

		builtBody := this.Body
		for parent != nil {
			contentRegexp := regexp.MustCompile("{% content %}")

			content := contentRegexp.ReplaceAllLiteral(parent.Body, builtBody)
			builtBody = content
			parent = layouts[parent.Parent]
		}

		// TODO: Regexp file extention should match only extensions defined in config
		includeRegexp := regexp.MustCompile("{% include .* %}")
		err := []string{}
		builtBody = includeRegexp.ReplaceAllFunc(builtBody, func(match []byte) []byte {
			nameRegexp := regexp.MustCompile(`[a-zA-Z0-9]*\.[a-zA-Z0-9]*`)
			includeName := nameRegexp.Find(match)

			if len(includeName) == 0 || includes[string(includeName)] == nil {
				err = append(err, "\""+string(includeName)+"\" invalid include name or include does not exist")
				return match
			} else if len(includes[string(includeName)].builtBody) == 0 {
				err = append(err, "\""+string(includeName)+"\" include has not been built")
				return match
			}
			return includes[string(includeName)].builtBody
		})

		if len(err) != 0 {
			return errors.New("could not build layouts\n\t" + strings.Join(err, "\n\t"))
		}

		layouts[name].builtBody = builtBody
	}

	return nil
}

func (l *Layout) showMeta() {
	fmt.Printf("Layout -> Name: %s Parent:%s Includes: [", l.Name, l.Parent)
	for _, dep := range l.Includes {
		fmt.Printf("%s ", dep)
	}
	fmt.Printf("]\n")
}

// * includes must have been read before this since it checks for valid includes
// * by definition a layout must have a 'content' tag
func NewLayout(name string, data []byte) (*Layout, error) {
	frontMatter, body, err := extractFrontMatterAndBody(data, "\n\n")
	if err != nil {
		return nil, err
	}

	var layout Layout = Layout{Name: name, Body: body}

	if len(frontMatter.Layout) != 0 {
		layout.Parent = frontMatter.Layout
	}

	contentRegexp := regexp.MustCompile("{% content %}")
	contentTags := contentRegexp.FindAll(body, -1)
	if len(contentTags) > 1 {
		err = errors.New(name + " layout can only contain one { % content } tag")
		return nil, err
	}

	includeRegexp := regexp.MustCompile("{% include .* %}")
	includeRegexp.ReplaceAllFunc(body, func(match []byte) []byte {
		nameRegexp := regexp.MustCompile(`[a-zA-Z0-9]*\.[a-zA-Z0-9]*`)
		includeName := nameRegexp.Find(match)
		if includeName == nil {
			err = errors.New(string(match) + " invalid include name")
		}

		layout.Includes = append(layout.Includes, string(includeName))
		return []byte{}
	})

	return &layout, err
}

func showLayouts(layouts map[string]*Layout, showRawBody bool, showBuiltBody bool) {
	for name, layout := range layouts {
		fmt.Printf("Layout: %s\n", name)

		var parent string
		if layout.Parent == "" {
			parent = "<none>"
		} else {
			parent = layout.Parent
		}
		fmt.Printf("Parent: %s\n", parent)

		deps := ""
		for _, includeName := range layout.Includes {
			deps = deps + includeName + ", "
		}
		if len(deps) == 0 {
			deps = "<none>"
		}
		fmt.Printf("Includes: %s\n", deps)

		if showRawBody {
			fmt.Printf("RawBody:\n%s\n", layout.Body)
		}

		if showBuiltBody {
			fmt.Printf("BuiltBody:\n%s\n", layout.builtBody)
		}
	}
}

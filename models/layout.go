package models

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

type Layouts map[string]Layout

// TODO: Add filename
type Layout struct {
	Name        string
	Body        []byte
	Status      int
	FrontMatter FrontMatter
	// More stuff needed here (each layout need access to .Posts and other site data)
	// each layout must contain 1 and only 1 {% content %} tag
}

/* path is absolute dir to folder of includes */
func (l Layouts) ReadDir(includes Includes, path string) (err error) {
	fmt.Printf("\nreading layouts from: %s\n", path)

	// TODO: Find better way to GLOB .html files

	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
		}
	}()

	for _, fileInfo := range fileInfos {
		filename := fileInfo.Name()
		ext := filepath.Ext(filename)

		if len(ext) == 0 || ext[1:] != "html" {
			continue
		}

		if l.has(filename) {
			continue
		}

		if err := l.buildLayout(includes, path, filename); err != nil {
			return err
		}
	}

	return err
}

// * includes must have been read before this since it checks for valid includes
// * by definition a layout must have a 'content' tag
func (l Layouts) buildLayout(includes Includes, path string, filename string) error {

	// if l.has(filename) {
	//     if (l.get(filename).Status & config.BUILT) > 0 {
	//         return nil
	//     } else {
	//         return errors.New("cyclic dependency detected in: " + filename)
	//     }
	// }

	uri := filepath.Join(path, filename)
	fmt.Println("parsing layout: " + filename)

	fileContent, err := ioutil.ReadFile(uri)
	if err != nil {
		return err
	}

	frontMatter, body, err := extractFrontMatterAndContent(fileContent)
	if err != nil {
		return err
	}

	layout := Layout{Body: body, Name: filename, FrontMatter: frontMatter}

	contentRegexp := regexp.MustCompile("{% content %}")

	_ = contentRegexp.FindAll(body, -1)

	// UNCOMMENT TO FORCE ONLY ONE CONTENT TAG
	// contentTags := contentRegexp.FindAll(body, -1)
	// if len(contentTags) > 1 {
	// 	return errors.New("layout may only contain one '{ % content }' tag")
	// }

	if frontMatter.HasLayout() {
		if l.has(frontMatter.Layout) {

			err = l.buildLayout(includes, path, frontMatter.Layout)
			if err != nil {
				panic("asdfasdf")
			}

			layout.Body = injectContentIntoLayout(layout.Body, l.get(frontMatter.Layout))
		}
	}

	if layout.Body, err = includes.ApplyIncludes(layout.Body); err != nil {
		return err
	}

	l[filename] = layout

	return err
}

func (self Layout) ApplyTo(data []byte) []byte {
	contentRegexp := regexp.MustCompile("{% content %}")
	return contentRegexp.ReplaceAllLiteral(self.Body, data)
}

func injectContentIntoLayout(content []byte, parent Layout) []byte {
	// debug: input
	// fmt.Printf("injecting: %s into parent: %s\n", layout.Name, parent.Name)
	// fmt.Printf("----------\n%s\n==========\n%s\n----------\n", string(layout.Body), string(parent.Body))

	contentRegexp := regexp.MustCompile("{% content %}")
	return contentRegexp.ReplaceAllLiteral(parent.Body, content)

	// debug: result
	// fmt.Println(string(layout.Body))
}

func (l Layouts) has(s string) bool {
	_, ok := l[s]
	return ok
}

func (l Layouts) get(s string) Layout {
	return l[s]
}

func (l Layouts) String() string {
	var buffer bytes.Buffer

	for _, layout := range l {
		s1 := fmt.Sprintf(
			`- Layout -------------------
            Name: %s
            %s`,
			layout.Name, layout.FrontMatter)
		s2 := fmt.Sprintf("%s\n", string(layout.Body))
		buffer.WriteString(s1)
		buffer.WriteString(s2)
	}

	return buffer.String()
}

// func applyLayout(layout string, body []byte, layouts map[string]*Layout) ([]byte, error) {
//     if len(layout) == 0 {
//         return body, nil
//     }

//     parentLayout := layouts[layout]
//     if parentLayout == nil {
//         return nil, errors.New("\"" + string(layout) + "\" layout does not exist")
//     }

//     contentRegexp := regexp.MustCompile("{% content %}")

//     return contentRegexp.ReplaceAllLiteral(parentLayout.Body, body), nil // TODO !! PROBABLY WRONG!
// }

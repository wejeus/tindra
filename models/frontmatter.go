package models

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v1"
	"regexp"
	"strings"
)

// TODO: Namechange to 'META' ?
type FrontMatter struct {
	// If set, this specifies the layout file to use.
	// Use the layout file name without the file extension.
	// Layout files must be placed in the  'layouts' directory.
	Layout string // TODO: Namechange LayoutName

	// Title is normaly generated using the filename if page is post.
	// If title is set in FrontMatter this title will be used instead.
	Title string

	// A date here overrides the date from the name of the post.
	// This can be used to ensure correct sorting of posts. Must have format YYYY-MM-DD.
	Date string // TODO

	// Similar to categories, one or multiple tags can be added to a post.
	// Also like categories, tags can be specified as a YAML list or a space- separated string.
	Tags []string // TODO

	// TODO
	// Instead of placing posts inside of folders, you can specify one or more
	// categories that the post belongs to. When the site is generated the post
	// will act as though it had been set with these categories normally.
	// Categories (plural key) can be specified as a YAML list or a space-separated string.
	// category
	// categories

	// TODO: Add support for custom tags?
}

func (f FrontMatter) HasLayout() bool {
	return len(f.Layout) != 0
}

func (f FrontMatter) String() string {
	return fmt.Sprintf("- FrontMatter --------------\n"+
		"Layout: %s\n"+
		"Title: %s\n"+
		"Date: %s\n"+
		"Published: %s\n"+
		"Tags: %s\n"+
		"----------------------------\n",
		f.Layout, f.Title, f.Date, f.Tags)
}

//  By default excerpt is your first paragraph of a post: everything before
//  the first two new lines. Testing if an excerpt is present is simply done by testing if first chars == "---"
//
//      ---
//      title: Example
//      ---
//
//      Second paragraph (post content)
//
// if file does not contain any frontMatter only the body will be returned and all other values will be nil
func extractFrontMatterAndContent(data []byte) (frontMatter FrontMatter, content []byte, err error) {
	separator := "\n\n"
	hasFrontMatter, err := regexp.Match("^---\n", data) // TODO: add better frontMatter regexp: "^---\n.*\n---\n"

	if hasFrontMatter {
		// if has frontMatter post must consist of 2 parts: (frontMatter, content)
		post := strings.SplitN(string(data), separator, 2)

		if len(post) != 2 {
			err = errors.New("could not extract (frontMatter, content)")
			return
		}

		// TODO: Check that frontMatter only contains one value for key 'layout'

		// TODO: Check that frontmatter contains date and date is valid

		err = yaml.Unmarshal([]byte(post[0]), &frontMatter)
		content = []byte(post[1])
	} else {
		content = data
	}

	return
}

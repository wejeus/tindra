package models

import (
    "errors"
    "fmt"
    "github.com/wejeus/tindra/utils"
    "gopkg.in/yaml.v1"
    "path/filepath"
    "regexp"
    "strings"
)

// TODO: Namechange to 'META' ?
type FrontMatter struct {
    // If set, this specifies the layout file to use.
    // Use the layout file name without the file extension.
    // Layout files must be placed in the  'layouts' directory.
    Layout string

    // Title is normaly generated using the filename if page is post.
    // If title is set in FrontMatter this title will be used instead.
    Title string

    // If you need your processed blog post URLs to be something other than
    // the default /year/month/day/title.html then you can set this variable
    // and it will be used as the final URL.
    Permalink string // TODO

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
        "Permalink: %s\n"+
        "Date: %s\n"+
        "Published: %s\n"+
        "Tags: %s\n"+
        "----------------------------\n",
        f.Layout, f.Title, f.Permalink, f.Date, f.Tags)
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
func extractFrontMatterAndContent(file string, data []byte, separator string) (frontMatter FrontMatter, content []byte, err error) {
    filename := filepath.Base(file)
    // TODO: Change regexp to use .MustCompile
    hasFrontMatter, err := regexp.Match("^---\n", data) // TODO: add better frontMatter regexp: "^---\n.*\n---\n"
    if hasFrontMatter {
        // if has frontMatter post must consist of 2 parts: (frontMatter, content)
        post := strings.SplitN(string(data), separator, 2)
        if len(post) != 2 {
            err = errors.New("could not extract (frontMatter, content)")
            return
        }

        // TODO: Check that frontMatter only contains one value for key 'layout'
        err = yaml.Unmarshal([]byte(post[0]), &frontMatter)
        content = []byte(post[1])
    } else {
        content = data
    }

    if len(frontMatter.Title) == 0 {
        frontMatter.Title = utils.TitleFromFilename(filename)
    }

    return
}

// func (p *Post) Date() string {
//     return "DDAAAATE"
// }

// if has permalink ->
//     return permalink

// var permalink

// if has date ->
//     permalink.append(date)
// else if filename has date ->
//     filenameDate = parseDate(filename)
//     permalink.append(filenameDate)

// permalink.append(filename)
// permalink.append(.html)

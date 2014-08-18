package context

import (
    "errors"
    "fmt"
    "gopkg.in/yaml.v1"
    "io/ioutil"
    "path/filepath"
    "regexp"
    "strings"
)

// TODO: Also read subfolders and prepend subfolder name to key
// TODO: Implement recursive read of dirs
// Reads all files with with postfix set to true in postfixes map.
func readFiles(directory string, postfixes map[string]bool) (files map[string][]byte, err error) {
    fmt.Printf("Reading directory: %s\n", directory)
    files = make(map[string][]byte)

    entries, err := ioutil.ReadDir(directory)
    if err != nil {
        return
    }

    for _, f := range entries {
        extension := filepath.Ext(f.Name())
        if len(extension) != 0 && postfixes[extension[1:]] {
            uri := filepath.Join(directory, f.Name())
            data, err := ioutil.ReadFile(uri)
            if err != nil {
                return files, err
            }
            files[f.Name()] = data
        }
    }

    return
}

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

    // Set to false if you don’t want a specific post to show up when the site is generated.
    // Defaults to true.
    Published string // TODO

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
func extractFrontMatterAndBody(data []byte, separator string) (frontMatter FrontMatter, body []byte, err error) {
    // TODO: Change regexp to use .MustCompile
    hasFrontMatter, err := regexp.Match("^---\n", data) // TODO: add better frontMatter regexp: "^---\n.*\n---\n"
    if hasFrontMatter {
        // if has frontMatter post must consist of 2 parts: (frontMatter, body)
        post := strings.SplitN(string(data), separator, 2)
        if len(post) != 2 {
            err = errors.New("could not extract (frontMatter, body)")
            return
        }

        // TODO: Check that frontMatter only contains one value for key 'layout'
        err = yaml.Unmarshal([]byte(post[0]), &frontMatter)
        body = []byte(post[1])
    } else {
        body = data
    }

    return
}

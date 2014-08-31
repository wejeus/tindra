package models

import (
    "bytes"
    //     "errors"
    "fmt"
    "github.com/russross/blackfriday"
    "io/ioutil"
    //     "log"
    "github.com/wejeus/tindra/config"
    "github.com/wejeus/tindra/utils"
    // "os"
    "path/filepath"
    // "regexp"
    "strings"
    //     "strings
)

type Posts map[string]Post

// Might share so many similarities with "Page" that we can merge them
type Post struct {
    Name        string
    Path        string // fullpath of file
    Template    string // template with html (layout/includes injected)
    Content     string // generated markdown
    Status      int
    FrontMatter FrontMatter

    // contains full path to all resources identified for post
    // will be put in /media possible under own subfolder
    Resources []Resource
}

// TODO: Change to something fancier (extract first paragraph and concat with title?)
func (self Post) Excerpt() string {
    head := self.Content[0 : len(self.Content)/5]
    return string(head)
}

// The 'posts' folder may contain .md documents in root or subfolders with .md files.
// Subfolder may contain resources. Subfolder may be 1 level deep.
func (posts Posts) ReadDir(layouts Layouts, path string) error {

    fileInfos, err := ioutil.ReadDir(path)
    if err != nil {
        return err
    }

    // add subdirs
    for _, file := range fileInfos {
        if file.IsDir() {

            subdirName := file.Name()
            subdirPath := filepath.Join(path, subdirName)

            markdownFiles, err := filepath.Glob(filepath.Join(subdirPath, "*.md"))
            if err != nil {
                fmt.Println("error: unknown error reading folder\nfolder: " + subdirPath)
                continue
            }

            if len(markdownFiles) != 1 {
                fmt.Println("error: post dir must contain one (and only one) markdown file \nfolder: " + subdirPath)
                continue
            }

            markdownFile := markdownFiles[0]
            markdownContent, err := ioutil.ReadFile(markdownFile)
            if err != nil {
                fmt.Println("error: could not read file \nfile: " + markdownFile)
                continue
            }

            resourceFiles, err := ioutil.ReadDir(subdirPath)
            if err != nil {
                fmt.Println("error: could not read folder\nfolder: " + subdirPath)
                continue
            }

            // read resources for dir
            var resourceFileList []Resource
            for _, resource := range resourceFiles {
                resourceFilenname := resource.Name()

                // TODO: Also ignore hidden files

                if filepath.Ext(resourceFilenname) != ".md" {
                    res := Resource{
                        Name:         resourceFilenname,
                        RelativeDir:  filepath.Join(config.RESOURCE_DIR_NAME, subdirName),
                        AbsolutePath: filepath.Join(subdirPath, resourceFilenname),
                    }

                    resourceFileList = append(resourceFileList, res)
                }
            }

            post, err := NewPost(layouts, markdownFile, markdownContent, resourceFileList)
            if err != nil {
                return err
            }

            posts[post.FrontMatter.Title] = post
        }
    }

    // add single .md files
    postsRootFolder := filepath.Join(path, "/*.md")
    matches, err := filepath.Glob(postsRootFolder)
    if err != nil {
        fmt.Println("error: unknown error reading folder\nfolder: " + postsRootFolder)
    }

    for _, file := range matches {
        markdownContent, err := ioutil.ReadFile(file)
        if err != nil {
            fmt.Println("error: could not read file \nfile: " + file)
            continue
        }

        post, err := NewPost(layouts, file, markdownContent, nil)
        if err != nil {
            return err
        }

        posts[post.FrontMatter.Title] = post
    }

    return nil
}

func NewPost(layouts Layouts, file string, content []byte, resources []Resource) (post Post, err error) {
    frontMatter, content, err := extractFrontMatterAndContent(file, content, "\n\n") // TODO: get from site config
    if err != nil {
        return
    }

    html := blackfriday.MarkdownCommon(content)
    body := injectContentIntoLayout(html, layouts.get(frontMatter.Layout))

    post = Post{
        Name:        filepath.Base(file),
        Path:        file,
        Template:    string(body),
        Content:     string(html),
        Status:      config.BUILT,
        FrontMatter: frontMatter,
        Resources:   resources,
    }

    return post, nil
}

// Extend to use different types of permalinks. For now uses format YYYY/MM/DD filename.html
func (self Post) Permalink() string {
    if len(self.FrontMatter.Permalink) != 0 {
        return self.FrontMatter.Permalink
    }

    date, filename := utils.SplitDateAndFilname(self.Name)
    if len(self.FrontMatter.Date) != 0 {
        date = self.FrontMatter.Date
    }

    if len(filename) == 0 {
        filename = self.Name
    }

    extension := filepath.Ext(filename)
    filename = filename[0 : len(filename)-len(extension)]

    filename = strings.Trim(filename, " _-")
    filename = filename + ".html" // TODO: Replace with 'addExtension-ish' lib function

    subpath := strings.Replace(date, "-", "/", 2)
    filename = strings.Replace(filename, " ", "_", -1)

    return filepath.Join(config.POSTS_DIR_NAME, subpath, filename)
}

func (posts Posts) String() string {
    var buffer bytes.Buffer

    for _, post := range posts {
        s1 := fmt.Sprintf("- Post ---------------------\n"+
            "Name: %s\n"+
            "Path: %s\n"+
            "Resources: %s\n"+
            "Status: %d\n"+
            "%s",
            post.Name, post.Path, post.Resources, post.Status, post.FrontMatter)
        s2 := fmt.Sprintf("= Template =================\n%s\n============================\n", string(post.Template))
        s3 := fmt.Sprintf("= Content ==================\n%s\n============================\n", string(post.Content))
        buffer.WriteString(s1)
        buffer.WriteString(s2)
        buffer.WriteString(s3)
    }

    return buffer.String()
}

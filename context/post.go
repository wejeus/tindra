package context

import (
    "bytes"
    "errors"
    "fmt"
    "github.com/russross/blackfriday"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "text/template"
)

type Post struct {
    Name        string // TODO: namechange PostTitle (or just Title everywhere?)
    Body        []byte
    Parent      string // nil if none
    FrontMatter FrontMatter
    Rendered    []byte
}

func readPostsDir(path string, allowedFiles map[string]bool) (posts map[string]*Post, err error) {
    path = filepath.Join(path, POSTS_DIR_NAME)

    posts = make(map[string]*Post)

    data, err := readFiles(path, allowedFiles)
    if err != nil {
        return
    }

    for name, raw := range data {
        p, postErr := NewPost(name, raw)
        if postErr != nil {
            err = postErr
            return
        }
        posts[name] = p
    }

    return
}

func NewPost(name string, data []byte) (post *Post, err error) {
    post = &Post{Name: name}

    frontMatter, body, err := extractFrontMatterAndBody(data, "\n\n") // TODO: get from site config
    if err != nil {
        return
    }

    var parent string = ""
    if len(frontMatter.Layout) != 0 {
        parent = frontMatter.Layout
    }

    post.Body = body
    post.Parent = parent
    post.FrontMatter = frontMatter

    return
}

// Assumes filename and path is absolute
func (p *Post) Build(site *Site) (err error) {
    fmt.Printf("building: %s\n", p.Name)

    html := blackfriday.MarkdownCommon(p.Body)
    rendered, err := applyLayout(p.Parent, html, site.layouts)
    fmt.Println(string(rendered))
    if err != nil {
        return err
    }

    page := TemplatePage{
        Site:      site,
        PageTitle: p.buildTitle(),
        Post:      *p,
        Date:      "asdfasdf",
        AllPosts:  site.posts,
    }

    // execute template
    t := template.Must(template.New("post").Parse(string(rendered)))
    var parsed bytes.Buffer
    if err := t.Execute(&parsed, page); err != nil {
        return err
    }
    p.Rendered = parsed.Bytes()

    // fmt.Println(parsed.String())
    return
}

func (p *Post) Install(site *Site, basePath string) (err error) {
    // build := p.buildInfo
    // if build == nil {
    // 	log.Fatal("need to build before install!")
    // }
    subpath, filename := p.Permalink()

    path := filepath.Join(basePath, BUILD_DIR_NAME, subpath)
    err = os.MkdirAll(path, os.ModeDir|0755)
    if err != nil {
        return
    }

    uri := filepath.Join(path, filename)

    fmt.Printf("installing: %s\n", uri)
    return ioutil.WriteFile(uri, p.Rendered, 0644)
}

func (p *Post) Date() string {
    return "DDAAAATE"
}

// TODO: Change to something fancier (extract first paragraph?)
func (p *Post) Excerpt() string {
    head := p.Body[0 : len(p.Body)/5]
    return string(blackfriday.MarkdownCommon(head))
}

// Extend to use different types of permalinks. For now uses format YYYY/MM/DD filename.html
func (p *Post) Permalink() (subpath, filename string) {
    date, filename, err := p.splitDateAndFilname()
    if err != nil {
        log.Fatal(err)
    }

    extension := filepath.Ext(filename)
    filename = filename[0 : len(filename)-len(extension)]

    filename = strings.Trim(filename, " _-")
    filename = filename + ".html" // TODO: Replace with 'addExtension-ish' lib function

    subpath = strings.Replace(date, "-", "/", 2)
    filename = strings.Replace(filename, " ", "_", -1)

    return
}

func (p *Post) buildTitle() string {
    _, name := p.Permalink()

    if p.FrontMatter.Title != "" {
        return p.FrontMatter.Title
    }

    // get 'title' part
    extension := filepath.Ext(name)
    title := name[0 : len(name)-len(extension)]

    // clean up title by removing possible whitespace or custom wordseparators
    title = strings.Replace(title, "-", " ", -1)
    title = strings.Replace(title, "_", " ", -1)
    title = strings.Trim(title, " ")

    return strings.Title(title)
}

// Post filename must be structured to include a filename prefix of "YYYY-MM-DD"
// What follows after the prefix is optional but it is encouraged to use the title for the post as name.
// Filename will be use as output to site so choose characters wisely.
// Returns error if date part of filname not valid.
func (p *Post) splitDateAndFilname() (date string, name string, err error) {
    bytes := []byte(p.Name)
    matcher := regexp.MustCompile("^[0-9]{4}-[0-9]{2}-[0-9]{2}.*")
    if !matcher.Match(bytes) {
        err = errors.New("post filename prefix does not match date convention")
        return
    }

    // extract date, title and replace extension
    date = string(bytes[:10])
    name = string(bytes[10:])

    return
}

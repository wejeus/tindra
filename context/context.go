package context

import (
    "bytes"
    "fmt"
    "github.com/wejeus/tindra/config"
    "github.com/wejeus/tindra/models"
    "github.com/wejeus/tindra/utils"
    "os"
    "path/filepath"
    "text/template" // TODO: Replate with "html/template" ?
)

// Site context, contains the constructs of the world
type Context struct {
    Includes models.Includes
    Layouts  map[string]models.Layout
    Posts    map[string]models.Post

    Resources []models.Resource // fullpath to resource dirs (/js, /media, /css, ..)
    // Pages    map[string]*Post // TODO
    // Resources []string // all the resouces we should copy to build dir (that is /js, /media, /css folders and maybe more?)
    // copy /favicon.ico
}

// Vill skilja på vad som är "inläst och halvbyggt" och vad som kommer genereras vid applicering av templates

func NewContext(path string) (context Context) {
    // config := NewDefaultConfig()
    // config.ReadFromConfigFile(path) // TODO: should be implicit in NewConfig()

    includes := models.Includes{}
    if err := includes.ReadDir(filepath.Join(path, config.INCLUDES_DIR_NAME)); err != nil {
        panic(err)
    }

    layouts := models.Layouts{}
    if err := layouts.ReadDir(includes, filepath.Join(path, config.LAYOUTS_DIR_NAME)); err != nil {
        panic(err)
    }

    posts := models.Posts{}
    if err := posts.ReadDir(layouts, filepath.Join(path, config.POSTS_DIR_NAME)); err != nil {
        panic(err)
    }
    // fmt.Printf("%s\n", posts)

    // TODO: Read pages (index.html, 404.html, others?)

    resources := []models.Resource{
        {"/media", "/", filepath.Join(path, "/media")},
        {"/css", "/", filepath.Join(path, "/css")},
        {"/js", "/", filepath.Join(path, "/js")},
        {"favicon.png", "/", filepath.Join(path, "/favicon.png")},
    }
    // TODO: Read resources: /js, /media, /css

    context = Context{Includes: includes, Layouts: layouts, Resources: resources, Posts: posts}

    return context
}

// func ReadResourceDir(path string) ([]models.Resource, error) {
//     var fileInfo os.FileInfo
//     if fileInfo, err := os.Stat(path); err != nil {
//         return nil, err
//     }

//     if !fileInfo.IsDir() {
//         return errors.New("given path is not a resource dir")
//     }

// }

// Installs to path
func (self Context) Install(path string) {
    fmt.Printf("\n")

    for postName, post := range self.Posts {
        fmt.Println("Installing: " + postName)
        // execute template using constructe
        construct := models.Construct{
            Posts:     self.Posts,
            SiteTitle: "SuperDuperTitle!",
            Page:      post,
        }

        finalHtmlPage := self.build(construct)
        outputPath := filepath.Join(path, post.Permalink())

        err := utils.WriteFile(outputPath, finalHtmlPage)
        if err != nil {
            panic(err)
        }
    }

    fmt.Printf("\n")
    for _, res := range self.Resources {
        err := utils.CopyFolderVerbatim(res.AbsolutePath, filepath.Join(path, res.RelativeDir))
        if err == nil {
            fmt.Printf("Copyied resource: %s\n", res.Name)
        }
    }

    // self.media.install()

    // self.javascripts.install()
    // self.css.install()
}

// must be used with a 'construct' since the construct is created after all post have been read and contains
// things such as pointers back and forth between posts, how many post in total and so on.
func (self Context) build(construct models.Construct) []byte {

    post := construct.Page

    // fmt.Print(string(post.Body))

    t, err := template.New("post").Parse(post.Template)
    if err != nil {
        panic(err)
    }

    t = template.Must(t, err)
    if err != nil {
        panic(err)
    }

    var parsed bytes.Buffer

    if err := t.Execute(&parsed, construct); err != nil {
        panic(err)
    }

    return parsed.Bytes()
}

// installs an array of paths to BUILD/media
//      js/file.js -> <root>/js/file.js
//      <something>/file.ext -> <root>/media/<something>/file.ext
func install(resources []models.Resource, buildPath string) error {

    for _, res := range resources {
        src := res.AbsolutePath
        dstDir := filepath.Join(buildPath, res.RelativeDir)
        dst := filepath.Join(buildPath, res.RelativeDir, res.Name)

        // make dir if not exists
        if err := os.MkdirAll(dstDir, os.ModeDir|0755); err != nil {
            panic(err)
        }

        fmt.Printf("%s -> %s\n", src, dst)
        err := utils.CopyFile(src, dst)
        if err != nil {
            panic(err)
        }
    }

    return nil
}

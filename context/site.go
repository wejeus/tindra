package context

import (
    "fmt"
    // "log"
    "path/filepath"
)

// Implement this on all 'page' types. Used when generating final site to get content
// type GetData interface{}
// TODO: Use go routines for templates "Once constructed, a template may be executed safely in parallel."

type Site struct {
    includes Includes
    layouts  map[string]*Layout
    posts    map[string]*Post
    // Pages    map[string]*Post // TODO
}

type RenderedSite struct {
    // Includes map[string][]byte
    // Includes map[string][]byte
}

// TODO: Better name
type TemplatePage struct {
    Site      *Site
    PageTitle string
    Post      Post // TODO should be *
    AllPosts  map[string]*Post
    Data      *Data
    Date      string // TODO
}

func NewSite(path string) (site *Site, err error) {
    fmt.Println("Generating new site...")

    config := NewDefaultConfig()
    config.ReadFromConfigFile(path) // TODO: should be implicit in NewConfig()

    // site = &Site{
    // 	basePath: path,
    // 	config:   config,
    // }

    // files, folders, err := getFilesAndFolders(config.basePath)

    // site.Data = NewData(config.prependAbsPath(DATA_DIR_NAME))
    // for k, v := range *site.Data {
    // 	fmt.Println(k)
    // 	for k, _ := range v {
    // 		fmt.Println("af " + k)
    // 	}
    // }

    site = new(Site)

    includes := Includes{}

    if err = includes.ReadDir(filepath.Join(path, INCLUDES_DIR_NAME)); err != nil {
        panic(err)
    }
    includes.printAllIncludes()

    // err = site.readLayoutsDir(path)
    // if err != nil {
    //     log.Fatal(err)
    // }

    // site.posts, err = readPostsDir(path, config.MarkdownExt)
    // if err != nil {
    // 	log.Fatal(err)
    // }

    // TODO: Read pages (index.html, 404.html, others?)
    return
}

func (site *Site) readLayoutsDir(path string) error {
    path = filepath.Join(path, LAYOUTS_DIR_NAME)

    allowedFiles := map[string]bool{"html": true}

    rawLayoutFiles, err := readFiles(path, allowedFiles)
    if err != nil {
        return err
    }

    layouts := make(map[string]*Layout)

    for name, raw := range rawLayoutFiles {
        l, layoutErr := NewLayout(name, raw)
        if layoutErr != nil {
            return layoutErr
        }
        layouts[name] = l
    }

    site.layouts = layouts

    return nil
}

// func (site *Site) Build() (*RenderedSite, error) {

//     builtSite := new(RenderedSite)

//     builtSite.Includes = make(map[string][]byte)

//     for name, rawInclude := range site.includes {
//         fmt.Println("Building: " + name)
//         rendered, err := rawInclude.Build(site)
//         if err != nil {
//             return nil, err
//         }

//         builtSite.Includes[name] = rendered
//     }

//     // if err := BuildLayouts(site.layouts, site.includes); err != nil {
//     // 	return err
//     // }

//     // for _, post := range site.posts {
//     // 	if err := post.Build(site); err != nil {
//     // 		return err
//     // 	}
//     // }

//     // TODO: We can build/generate CSS here

//     return builtSite, nil
// }

func (site *RenderedSite) Install() error {
    fmt.Println("Installing..")
    // copy /css
    // copy /javascript
    // copy /images
    // copy /favicon.ico

    // output previously built posts
    // output previously built pages (index.html, 404.html, others?)

    // write Posts to a dir and filename based on permalink/title

    // buildPath := filepath.Join(site.basePath, BUILD_DIR_NAME)
    // fmt.Println("build: " + buildPath)
    // if err := os.RemoveAll(buildPath); err != nil {
    // 	return err
    // }

    // if err := os.MkdirAll(buildPath, os.ModeDir|0755); err != nil {
    // 	return err
    // }

    // if err := copyFolderVerbatim(filepath.Join(site.basePath, CSS_DIR_NAME), buildPath); err != nil {
    // 	return err
    // }

    // for _, post := range site.posts {
    // 	if err := post.Install(site); err != nil {
    // 		return err
    // 	}
    // }

    return nil
}

// TODO
// func clean()

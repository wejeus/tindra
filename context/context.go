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
	BasePath string
	Includes models.Includes
	Layouts  map[string]models.Layout

	Posts     map[string]*models.Post
	Resources []models.Resource // fullpath to resource dirs (/js, /media, /css, ..)
	Pages     map[string]*models.Page
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

	pages := models.NewPages(includes, layouts, filepath.Join(path, config.PAGES_DIR_NAME))
	// fmt.Printf("%s\n", pages)

	resources := []models.Resource{
		{config.RESOURCE_DIR_NAME, "/"},
		{config.CSS_DIR_NAME, "/"},
		{config.JS_DIR_NAME, "/"},
		{"favicon.png", "/"},
	}
	// TODO: Read resources: /js, /media, /css

	context = Context{BasePath: path, Includes: includes, Layouts: layouts, Resources: resources, Posts: posts, Pages: pages}

	return context
}

func InstallResources(resources []models.Resource, basePath, buildPath string) {
	for _, res := range resources {
		in := filepath.Join(basePath, res.RelPath())

		// fmt.Printf("Installing resource: %s\n", res.RelPath())
		fileInfo, err := os.Stat(in)

		if err == nil && !fileInfo.IsDir() {
			out := filepath.Join(buildPath, res.Filename)
			fmt.Printf("Installing resource: %s\n", out)
			utils.CopyFile(in, out)
		} else {
			out := buildPath
			fmt.Printf("Installing resource: %s\n", out)
			utils.CopyFolderRecursive(in, out)
		}
	}
}

func (self Context) InstallPosts(buildPath string) {
	for _, post := range self.Posts {
		outputPath := filepath.Join(buildPath, post.Permalink())
		fmt.Println("Installing post: " + outputPath)

		construct := models.Construct{
			Posts:     self.Posts,
			SiteTitle: "SuperDuperTitle!",
			Page:      post,
		}

		html := self.build(construct)
		utils.WriteDataToFile(outputPath, html)

		InstallResources(post.Resources, self.BasePath, filepath.Dir(outputPath))
	}
}

func (self Context) InstallPages(buildPath string) {
	for pageName, page := range self.Pages {
		outputPath := filepath.Join(buildPath, pageName)
		fmt.Println("Installing page: " + outputPath)

		construct := models.Construct{
			Posts:     self.Posts,
			SiteTitle: "SuperDuperTitle!",
		}

		html := construct.Build(string(page.Template))
		utils.WriteDataToFile(outputPath, html)
	}
}

// Installs to path
func (self Context) Install(buildPath string) {
	InstallResources(self.Resources, self.BasePath, buildPath)
	self.InstallPages(buildPath)
	self.InstallPosts(buildPath)
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

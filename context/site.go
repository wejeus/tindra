package context

import (
	"log"
	"os"
	// "path/filepath"
	"fmt"
	"io/ioutil"
)

// Implement this on all 'page' types. Used when generating final site to get content
// type GetData interface{}
// TODO: Use go routines for templates "Once constructed, a template may be executed safely in parallel."

type Site struct {
	config   *Config
	Data     *Data
	Includes map[string]*Include
	Layouts  map[string]*Layout
	// Pages    map[string]*Post // TODO
	Posts map[string]*Post
}

func NewSite() (site *Site, err error) {
	fmt.Print("Generating new site...")

	config := NewConfig()
	config.ReadFromConfigFile() // TODO: should be implicit in NewConfig()

	site = &Site{
		config: config,
	}

	// files, folders, err := getFilesAndFolders(config.basePath)

	site.Data = NewData(config.prependAbsPath(DATA_DIR_NAME))

	// for k, v := range *site.Data {
	// 	fmt.Println(k)
	// 	for k, _ := range v {
	// 		fmt.Println("af " + k)
	// 	}
	// }

	site.Includes, err = readIncludesDir(config.prependAbsPath(INCLUDES_DIR_NAME))
	if err != nil {
		log.Fatal(err)
	}
	err = BuildIncludes(site.Includes)
	if err != nil {
		log.Fatal(err)
	}

	site.Layouts, err = readLayoutsDir(config.prependAbsPath(LAYOUTS_DIR_NAME))
	if err != nil {
		log.Fatal(err)
	}
	err = BuildLayouts(site.Layouts, site.Includes)
	if err != nil {
		log.Fatal(err)
	}

	site.Posts, err = readPostsDir(config.prependAbsPath(POSTS_DIR_NAME), config.MarkdownExt)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range site.Posts {
		err = post.Build(site)
		if err != nil {
			log.Fatal(err)
		}
	}

	// for _, post := range site.Posts {
	// 	fmt.Println(string(post.Rendered))
	// }
	// site.Posts, err = site.loadAndPreprocessPosts(config.prependAbsPath(POSTS_DIR_NAME), config.MarkdownExt)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// s.Pages, err = s.loadAndPreprocessPages(config.basePath, map[string]bool{"html": true})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	return
}

func copyFolder(src, dest string) {

}

// separates files and folders, skips special entries and hidden entries
func getFilesAndFolders(path string) (specialFiles, verbatimDirs map[string]bool, err error) {
	special := map[string]bool{
		MAIN_CONFIG_FILENAME: true,
		INCLUDES_DIR_NAME:    true,
		LAYOUTS_DIR_NAME:     true,
		POSTS_DIR_NAME:       true,
		DATA_DIR_NAME:        true,
		PLUGINS_DIR_NAME:     true,
		BUILD_DIR_NAME:       true,
	}

	specialFiles = make(map[string]bool)
	verbatimDirs = make(map[string]bool)

	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if special[entry.Name()] || string([]byte(entry.Name())[0:1]) == "." {
			continue
		}

		if entry.IsDir() {
			verbatimDirs[entry.Name()] = true
		} else {
			specialFiles[entry.Name()] = true
		}
	}

	return
}

func (site *Site) BuildAndInstall() (err error) {
	// FIXME: Dont dare to run this, maybe it removes everything!
	// err = os.RemoveAll(site.config.BuildPath)
	// if err != nil {
	// 	return
	// }

	err = os.MkdirAll(site.config.getAbsBuildPath(), os.ModeDir|0755)
	if err != nil {
		return
	}

	// for _, p := range site.Pages {
	// 	err = p.buildAndInstall(site)
	// 	if err != nil {
	// 		return
	// 	}
	// }

	return
}

func list(m map[string][]byte) {
	for name, _ := range m {
		fmt.Printf("%s ", name)
	}
}

func show(m map[string][]byte) {
	for name, data := range m {
		fmt.Printf("%s\n----------------------------\n%s\n", name, string(data))
	}
}

// func installDir(directory string) {
// 	os.Wal
// 	// Write to disk
// 	err = os.MkdirAll(outPath, os.ModeDir|0755)
// 	if err != nil {
// 		return
// 	}
// 	fmt.Printf("installing: %s\n", outFile)
// 	return ioutil.WriteFile(outFile, parsed.Bytes(), 0644)
// }

// func (site *Site) loadAndPreprocessPosts(directory string, postfixes map[string]bool) (posts map[string]*Post, err error) {
// 	files, err := readFiles(directory, postfixes)
// 	if err != nil {
// 		return
// 	}

// 	posts = make(map[string]*Post)
// 	for filename, data := range files {
// 		post := site.NewPost(filename, data)
// 		posts[filename] = post
// 	}

// 	return
// }

// func (site *Site) loadAndPreprocessStaticPages(directory string, postfixes map[string]bool) (posts map[string]*Post, err error) {
// 	files, err := readFiles(directory, postfixes)
// 	if err != nil {
// 		return
// 	}

// 	posts = make(map[string]*Post)
// 	for filename, data := range files {
// 		fmt.Print("parsing page: " + filename)
// 		post := site.NewPage(filename, data)
// 		posts[filename] = post
// 	}

// 	return
// }

// func (site *Site) NewStaticPage(filename string, data []byte) (post *Post) {
// 	frontMatter, body, err := extractFrontMatterAndBody(data, "\n\n") // TODO: get from site config
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if frontMatter.Layout != "" && site.Layouts[frontMatter.Layout] == nil {
// 		log.Fatal("error: \"" + filename + "\" declared layout that do not exist!")
// 	}

// 	title := site.config.Name
// 	if len(frontMatter.Title) != 0 {
// 		title = frontMatter.Title
// 	}

// 	parentLayout := site.Layouts[frontMatter.Layout]

// 	return &Post{
// 		parent:   parentLayout,
// 		content:  body,
// 		filename: filename,
// 		path:     ".",

// 		Title: title,
// 	} // TODO: Does this cause a copy upon return?
// }

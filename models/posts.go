package models

import (
	"bytes"
	"fmt"
	"github.com/wejeus/tindra/config"
	"io/ioutil"
	"path/filepath"
)

type Posts map[string]*Post

// The 'posts' folder may contain .md documents in root or subfolders with .md files.
// Subfolder may contain resources. Subfolder may be 1 level deep.
func (posts Posts) ReadDir(layouts Layouts, path string) error {

	// TODO: results in runtime panic if "posts" folder does not exists!

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

				// TODO: ignore hidden files and others marked as DO_NOT_COPY

				res := Resource{
					Filename: resourceFilenname,
					RelDir:   filepath.Join(config.POSTS_DIR_NAME, subdirName),
				}

				resourceFileList = append(resourceFileList, res)
			}

			post := NewPost(path, markdownFile, markdownContent, layouts, resourceFileList)

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
			fmt.Println("could not read file \nfile: " + file)
			continue
		}

		post := NewPost(path, file, markdownContent, layouts, nil)

		posts[post.FrontMatter.Title] = post
	}

	return nil
}

func (self Posts) String() string {
	var buffer bytes.Buffer
	for _, post := range self {
		buffer.WriteString(post.String())
	}
	return buffer.String()
}

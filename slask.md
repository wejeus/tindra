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
package models

import (
	"bytes"
	"fmt"
	"github.com/russross/blackfriday"
	"github.com/wejeus/tindra/config"
	"path/filepath"
	"strings"
)

// Rel: posts/droidcon
// Name: my_awesome_talk.md (check: must be valid format, only one .md file)
// Resources: [my_awesome_talk.md, image1.png, subdir/image2.png]

// <required>
// Date
// Title
// Layout

// <generated>
// Permalink: posts/2014/05/27/my_awesome_talk.html
// Output path: posts/2014/05/27/

// When building: give warning on duplicate resources

// Might share so many similarities with "Page" that we can merge them
// Path        string // fullpath of file
type Post struct {
	Filename    string     // ex: my_awesome_post.md
	RelDir      string     // ex: posts/droidcon/
	Resources   []Resource // contains full path to all resources identified for post
	FrontMatter FrontMatter

	Template string // template with html (layout/includes injected)
	Content  string // generated markdown
}

func NewPost(base, file string, content []byte, layouts Layouts, resources []Resource) *Post {
	frontMatter, content, err := extractFrontMatterAndContent(content)
	if err != nil {
		panic("[post.NewPost] missing or incorrect FrontMatter. Post " + file)
	}

	html := blackfriday.MarkdownCommon(content)

	if !layouts.has(frontMatter.Layout) {
		panic("[post.NewPost] declared layout does not exists. Missing layout " + frontMatter.Layout + " in post " + file)
	}

	body := injectContentIntoLayout(html, layouts.get(frontMatter.Layout))
	rel, _ := filepath.Rel(base, filepath.Dir(file))
	post := Post{
		Filename:    filepath.Base(file),
		RelDir:      filepath.Join(config.POSTS_DIR_NAME, rel),
		Template:    string(body),
		Content:     string(html),
		FrontMatter: frontMatter,
		Resources:   resources,
	}

	return &post
}

// unique idenifier (relative path + filename)
// ex: droidcon/my_awesome_talk.md
func (self Post) RelPath() string {
	return filepath.Join(self.RelDir, self.Filename)
}

// TODO: Change to something fancier (extract first paragraph and concat with title?)
func (self Post) Excerpt() string {
	head := self.Content[0 : len(self.Content)/5]
	return string(head)
}

func NameWithExtension(name, ext string) string {
	currentExt := filepath.Ext(name)

	if len(currentExt) > 0 {
		name = name[0 : len(name)-len(currentExt)]
	}

	return name + ext
}

// POSSIBLE TODO: Extend to use different types of permalinks. For now uses format YYYY/MM/DD filename.html
// generates: posts/2014/05/27/my_awesome_talk.html
func (self Post) Permalink() string {
	datepath := strings.Replace(self.FrontMatter.Date, "-", "/", 2)
	noWhiteSpaceName := strings.Replace(self.Filename, " ", "_", -1)
	htmlFilename := NameWithExtension(noWhiteSpaceName, ".html")

	return filepath.Join(config.POSTS_DIR_NAME, datepath, htmlFilename)
}

func (self Post) String() string {
	var buffer bytes.Buffer

	s1 := fmt.Sprintf("- Post ---------------------\n"+
		"Filename: %s\n"+
		"RelDir: %s\n"+
		"Resources: \n%s\n\n"+
		"RelPath(): %s\n"+
		"Permalink(): %s\n"+
		"%s",
		self.Filename, self.RelDir, self.Resources, self.RelPath(), self.Permalink(), self.FrontMatter)
	s2 := fmt.Sprintf("= Template =================\n%s\n============================\n", string(self.Template))
	s3 := fmt.Sprintf("= Content ==================\n%s\n============================\n", string(self.Content))
	buffer.WriteString(s1)
	buffer.WriteString(s2)
	buffer.WriteString(s3)

	return buffer.String()
}

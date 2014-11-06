package models

import (
	"bytes"
	"text/template" // TODO: Replate with "html/template" ?
)

// TODO: Separet constructs for pages/posts?
type Construct struct {
	WebRoot   string
	Posts     map[string]*Post
	SiteTitle string
	Page      *Post // TODO: Refactor to methods so we can add some error handling
	// Pages    map[string]*Post // TODO
	// Resources []string // all the resouces we should copy to build dir (that is /js, /media, /css folders and maybe more?)
}

func (self Construct) Build(src string) []byte {
	t, err := template.New("page").Parse(src)
	if err != nil {
		panic(err)
	}

	t = template.Must(t, err)
	if err != nil {
		panic(err)
	}

	var parsed bytes.Buffer

	if err := t.Execute(&parsed, self); err != nil {
		panic(err)
	}

	return parsed.Bytes()
}

func (self Construct) PreviousPost() string {
	// not implemented yet
	return "permalink:previous"
}

func (self Construct) NextPost() string {
	// not implemented yet
	return "permalink:next"
}

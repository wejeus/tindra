package models

// TODO: Separet constructs for pages/posts?
type Construct struct {
    Posts     map[string]Post
    SiteTitle string
    Page      Post // TODO: Refactor to methods so we can add some error handling
    // Pages    map[string]*Post // TODO
    // Resources []string // all the resouces we should copy to build dir (that is /js, /media, /css folders and maybe more?)
}

func (self Construct) PreviousPost() string {
    // not implemented yet
    return "permalink:previous"
}

func (self Construct) NextPost() string {
    // not implemented yet
    return "permalink:next"
}

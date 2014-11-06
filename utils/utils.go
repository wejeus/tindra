package utils

import (
	"path/filepath"
	"regexp"
	"strings"
)

func TitleFromFilename(filename string) string {
	// get 'title' part
	extension := filepath.Ext(filename)
	title := filename[0 : len(filename)-len(extension)]

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
func SplitDateAndFilname(filename string) (date, name string) {
	bytes := []byte(filename)
	matcher := regexp.MustCompile("^[0-9]{4}-[0-9]{2}-[0-9]{2}.*")
	if !matcher.Match(bytes) {
		return
	}

	// extract date, title and replace extension
	date = string(bytes[:10])
	name = string(bytes[10:])

	return
}

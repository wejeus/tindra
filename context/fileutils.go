package context

import (
	"path/filepath"
)

func (site *Site) getPath(dir string) {
	return filepath.Join(site.config.basePath, dir)
}

func (site *Site) getBuildPath(dir string) {
	return filepath.Join(site.config.basePath, site.config.BuildPath, dir)
}

// TODO: Add possible regex match on filename instead?
// TODO: Also read subfolders and prepend subfolder name to key
// TODO: function name change: readFilesRaw?
// TODO: Implement recursive read of dirs
// Reads all files with with postfix set to true in postfixes map.
// A new map of {parentPaths/filename, content} is returned or error if something goes bananas.
func readFiles(directory string, postfixes map[string]bool) (files map[string][]byte, err error) {
	log.Printf("Reading directory: %s\n", directory)
	files = make(map[string][]byte)

	entries, err := ioutil.ReadDir(directory)
	if err != nil {
		return
	}

	for _, f := range entries {
		extension := filepath.Ext(f.Name())
		if len(extension) != 0 && postfixes[extension[1:]] {
			uri := filepath.Join(directory, f.Name())
			data, err := ioutil.ReadFile(uri)
			if err != nil {
				return m, err
			}
			files[f.Name()] = data
		}
	}

	return
}

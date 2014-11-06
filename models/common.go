package models

import (
	"fmt"
	"path/filepath"
)

type Resource struct {
	Filename string
	RelDir   string
}

func (self Resource) RelPath() string {
	return filepath.Join(self.RelDir, self.Filename)
}

func (self Resource) String() string {
	return fmt.Sprintf("\tFilename: %s\nRelDir: %s\nRelPath(): %s\n", self.Filename, self.RelDir, self.RelPath())
}

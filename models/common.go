package models

import "fmt"

type Resource struct {
    Name         string
    RelativeDir  string
    AbsolutePath string
}

func (self Resource) String() string {
    return fmt.Sprintf("Name: %s\nRelativeDir: %s\nAbsolutePath: %s\n", self.Name, self.RelativeDir, self.AbsolutePath)
}

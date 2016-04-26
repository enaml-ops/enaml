package pkg

import (
	"archive/zip"
	"fmt"
	"regexp"
)

// zipWalker walks a .zip (or .pivotal) file implementing the Walker interface
type zipWalker struct {
	zipPath   string
	callbacks map[*regexp.Regexp]WalkFunc
}

// NewZipWalker creates a new Walker instance that can read a .zip stream
func NewZipWalker(zipFile string) Walker {
	return zipWalker{
		zipPath:   zipFile,
		callbacks: make(map[*regexp.Regexp]WalkFunc),
	}
}

func (z zipWalker) OnMatch(regexExpr string, fn WalkFunc) {
	regex, err := regexp.Compile(regexExpr)
	if err != nil {
		// this indicates a programming error, not something that should normally happen
		panic(err.Error())
	}
	z.callbacks[regex] = fn
}

func (z zipWalker) Walk() error {
	zr, err := zip.OpenReader(z.zipPath)
	if err != nil {
		return err
	}
	for _, zipFile := range zr.File {
		if !zipFile.FileInfo().IsDir() {
			fmt.Println(zipFile.Name)
			for regex, fn := range z.callbacks {
				if regex.MatchString(zipFile.Name) {
					r, err := zipFile.Open()
					if err != nil {
						return err
					}
					fn(FileEntry{
						FileName: zipFile.Name,
						Reader:   r,
					})
				}
			}
		}
	}
	return nil
}

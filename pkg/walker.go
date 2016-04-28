// Package pkg provides utilties for reading from compressed archives.
package pkg

import "io"

// FileEntry is a package entry for a file.
type FileEntry struct {
	FileName string
	Reader   io.Reader
}

// WalkFunc is the function signature for file registration callbacks
// while walking a package.
type WalkFunc func(f FileEntry) error

// Walker walks a package invoking a callback for each matching file entry.
type Walker interface {
	OnMatch(regexExpr string, fn WalkFunc)
	Walk() error
}

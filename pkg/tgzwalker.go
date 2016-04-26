package pkg

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"regexp"
)

// TgzWalker walks a .tgz file implementing the Walker interface
type tgzWalker struct {
	pkgReader io.Reader
	callbacks map[*regexp.Regexp]WalkFunc
}

// NewTgzWalker creates a new Walker instance that can read a .tgz stream
func NewTgzWalker(pkgReader io.Reader) Walker {
	return tgzWalker{
		pkgReader: pkgReader,
		callbacks: make(map[*regexp.Regexp]WalkFunc),
	}
}

func (t tgzWalker) OnMatch(regexExpr string, fn WalkFunc) {
	regex, err := regexp.Compile(regexExpr)
	if err != nil {
		// this indicates a programming error, not something that should normally happen
		panic(err.Error())
	}
	t.callbacks[regex] = fn
}

func (t tgzWalker) Walk() error {
	gr, err := gzip.NewReader(t.pkgReader)
	if err != nil {
		return err
	}
	tr := tar.NewReader(gr)
	for {
		h, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if h.Typeflag == tar.TypeReg {
			for regex, fn := range t.callbacks {
				if regex.MatchString(h.Name) {
					fn(FileEntry{
						FileName: h.Name,
						Reader:   tr,
					})
				}
			}
		}
	}
	return nil
}

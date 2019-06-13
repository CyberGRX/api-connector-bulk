package assets

import (
	"net/http"
	"strings"
)

type EmbeddedFileSystem struct {
	fs http.FileSystem
}

func (b *EmbeddedFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *EmbeddedFileSystem) Exists(prefix string, filepath string) bool {
	if _, err := b.fs.Open(filepath[1:]); err == nil {
		return true
	}

	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func NewEmbeddedFileSystem() *EmbeddedFileSystem {
	return &EmbeddedFileSystem{fs: Assets}
}

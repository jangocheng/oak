package file

import (
	"io"
	"os"
)

type File interface {
	io.ReadCloser
	Stat() (os.FileInfo, error)
}
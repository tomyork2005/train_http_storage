package storages

import "errors"

var (
	ErrFileNotFound = errors.New("File not found")
)

type FileToAdd struct {
	Alias      string
	PathToFile string
	UserId     int64
}

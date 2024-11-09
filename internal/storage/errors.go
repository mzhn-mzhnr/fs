package storage

import "errors"

var (
	ErrFileAlreadyExists = errors.New("file already exists")
	ErrFileNotExists     = errors.New("file not exists")
	ErrFileDuplicate     = errors.New("file duplicate")
	ErrFileBadId         = errors.New("file bad id")
)

package domain

import (
	"bytes"
	"io"
	"net/http"
)

type Entry struct {
	filename string
	io.Reader
}

func NewEntry(filename string, reader io.Reader) *Entry {
	return &Entry{filename: filename, Reader: reader}
}

func (e *Entry) Filename() string {
	return e.filename
}

func (e *Entry) ContentType() (string, error) {

	buf := new(bytes.Buffer)

	if _, err := io.Copy(buf, e); err != nil {
		return "", err
	}

	ct := http.DetectContentType(buf.Bytes())

	e.Reader = bytes.NewReader(buf.Bytes())
	return ct, nil
}

type FileRecord struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

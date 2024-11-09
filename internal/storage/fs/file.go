package fs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mzhn/fileservice/internal/domain"
	"mzhn/fileservice/internal/storage"
	"mzhn/fileservice/pkg/sl"
	"os"
	"time"
)

// File implements fileservice.FileProvider.
func (f *FileStorage) File(ctx context.Context, filename string) (*domain.Entry, error) {
	fn := "fs.File"
	log := f.logger.With(sl.Method(fn), slog.String("filename", filename))

	path := f.joinFileName(filename)

	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%s: %w", fn, storage.ErrFileNotExists)
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	log.Debug("start file read")
	start := time.Now()
	if _, err := io.Copy(buf, file); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	log.Debug("file read", sl.Millis("took", time.Since(start)))

	return domain.NewEntry(filename, buf), nil
}

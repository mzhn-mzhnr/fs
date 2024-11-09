package fs

import (
	"context"
	"fmt"
	"log/slog"
	"mzhn/fileservice/internal/domain"
	"mzhn/fileservice/pkg/sl"
	"strings"
	"sync"
)

// UploadMany implements fileservice.FileUploader.
func (f *FileStorage) UploadMany(ctx context.Context, entries []*domain.Entry) (err error) {

	fn := "fs.UploadMany"
	log := f.logger.With(sl.Method(fn))

	wg := &sync.WaitGroup{}

	type written struct {
		filename string
		err      error
	}

	done := make(chan written)
	failed := make([]string, 0, len(entries))

	wg.Add(len(entries))

	for _, entry := range entries {
		go func(e *domain.Entry) {
			defer wg.Done()
			done <- written{
				filename: entry.Filename(),
				err:      f.upload(ctx, entry),
			}
		}(entry)
	}

	go func() {
		for w := range done {
			if w.err != nil {
				log.Error("file upload failed", slog.String("filename", w.filename), sl.Err(w.err))
				failed = append(failed, fmt.Sprintf("%s: %s", w.filename, w.err))
			}
		}
	}()

	wg.Wait()
	close(done)

	if len(failed) != 0 {
		return fmt.Errorf("%s: %s", fn, strings.Join(failed, ", "))
	}

	return nil
}

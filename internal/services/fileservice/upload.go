package fileservice

import (
	"context"
	"fmt"
	"log/slog"
	"mzhn/fileservice/internal/domain"
	"mzhn/fileservice/pkg/sl"
)

func (s *FileService) Upload(ctx context.Context, entry *domain.Entry) (string, error) {

	fn := "fileservice.Upload"

	id := s.idProvider.Provide()
	filename := entry.Filename()

	log := s.logger.With(sl.Method("fn"), slog.String("filename", filename), slog.String("id", id))

	log.Info("uploading a file")
	if err := s.uploader.Upload(ctx, domain.NewEntry(id, entry)); err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	log.Info("saving file")
	if err := s.saver.Save(ctx, domain.FileRecord{Id: id, Name: filename}); err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return formatPath(id, filename), nil
}

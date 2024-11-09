package fileservice

import (
	"context"
	"fmt"
	"log/slog"
	"mzhn/fileservice/internal/domain"
	"mzhn/fileservice/pkg/sl"
)

func (s *FileService) Upload(ctx context.Context, entry *domain.Entry) (*domain.FileRecord, error) {

	fn := "fileservice.Upload"

	id := s.idProvider.Provide()
	filename := entry.Filename()

	log := s.logger.With(sl.Method("fn"), slog.String("filename", filename), slog.String("id", id))

	log.Info("uploading a file")
	if err := s.uploader.Upload(ctx, domain.NewEntry(id, entry)); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	log.Info("saving file")
	rec := &domain.FileRecord{Id: id, Name: filename}
	if err := s.saver.Save(ctx, rec); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return rec, nil
}

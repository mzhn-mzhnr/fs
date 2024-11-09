package fileservice

import (
	"context"
	"mzhn/fileservice/internal/domain"
	"mzhn/fileservice/pkg/sl"

	"github.com/samber/lo"
)

func (s *FileService) UploadMany(ctx context.Context, entries []*domain.Entry) ([]*domain.FileRecord, error) {

	fn := "fileservice.UploadMany"
	log := s.logger.With(sl.Method(fn))

	log.Info("uploading many files")

	ee := lo.Map(entries, func(e *domain.Entry, i int) *domain.Entry {
		return domain.NewEntry(s.idProvider.Provide(), e)
	})

	if err := s.uploader.UploadMany(ctx, ee); err != nil {
		return nil, err
	}

	records := lo.Map(entries, func(e *domain.Entry, i int) *domain.FileRecord {
		return &domain.FileRecord{Id: ee[i].Filename(), Name: e.Filename()}
	})

	if err := s.saver.SaveMany(ctx, records); err != nil {
		return nil, err
	}

	return records, nil
}

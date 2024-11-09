package fileservice

import (
	"context"
	"fmt"
	"log/slog"
	"mzhn/fileservice/internal/domain"
	"mzhn/fileservice/pkg/sl"
)

type FileUploader interface {
	Upload(ctx context.Context, entry *domain.Entry) error
	UploadMany(ctx context.Context, entries []*domain.Entry) error
}

type FileDeleter interface {
	Delete(ctx context.Context, filename string) error
}

type FileProvider interface {
	File(ctx context.Context, filename string) (*domain.Entry, error)
}

type IdProvider interface {
	Provide() string
}

type FileSaver interface {
	Save(ctx context.Context, record *domain.FileRecord) error
	SaveMany(ctx context.Context, records []*domain.FileRecord) error
}

type FileService struct {
	uploader   FileUploader
	provider   FileProvider
	saver      FileSaver
	idProvider IdProvider

	logger *slog.Logger
}

func New(uploader FileUploader, provider FileProvider, saver FileSaver, ip IdProvider) *FileService {
	return &FileService{uploader, provider, saver, ip, slog.With(sl.Module("FileService"))}
}

func formatPath(id, name string) string {
	return fmt.Sprintf("/file/%s?id=%s", name, id)
}

func (s *FileService) File(ctx context.Context, id string) (*domain.Entry, error) {
	return s.provider.File(ctx, id)
}

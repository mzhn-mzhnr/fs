package fs

import (
	"errors"
	"fmt"
	"log/slog"
	"mzhn/fileservice/internal/config"
	"mzhn/fileservice/internal/services/fileservice"
	"mzhn/fileservice/pkg/sl"
	"os"
)

var _ fileservice.FileProvider = (*FileStorage)(nil)
var _ fileservice.FileUploader = (*FileStorage)(nil)

type FileStorage struct {
	staticPath string
	logger     *slog.Logger
}

func (f *FileStorage) createVolume() error {
	return os.Mkdir(f.staticPath, 0755)
}

func (f *FileStorage) checkVolume() error {
	fn := "fs.checkVolume"
	log := f.logger.With(sl.Method(fn))

	stat, err := os.Stat(f.staticPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return f.createVolume()
		}
	}

	if !stat.IsDir() {
		log.Error("current path to static points to file not a directory", slog.String("staticPath", f.staticPath))
		panic(fmt.Errorf("current path to static points to file not a directory"))
	}

	return nil
}

func (f *FileStorage) joinFileName(filename string) string {
	return fmt.Sprintf("%s/%s", f.staticPath, filename)
}

func New(cfg *config.Config) *FileStorage {
	fs := &FileStorage{
		staticPath: cfg.FS.Path,
		logger:     slog.With(sl.Module("fs.FileStorage")),
	}

	_ = fs.checkVolume()

	return fs
}

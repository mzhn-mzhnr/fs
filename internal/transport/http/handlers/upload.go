package handlers

import (
	"log/slog"
	"mzhn/fileservice/internal/domain"
	"mzhn/fileservice/internal/services/fileservice"
	"mzhn/fileservice/pkg/sl"

	"github.com/labstack/echo/v4"
)

func Upload(svc *fileservice.FileService) echo.HandlerFunc {
	return func(c echo.Context) error {
		form, err := c.MultipartForm()
		if err != nil {
			slog.Error("failed to get multipart form", sl.Err(err))
			return c.JSON(echo.ErrInternalServerError.Code, throw("internal server error"))
		}

		files := form.File["file"]

		ctx := c.Request().Context()

		paths := make([]string, 0, len(files))

		if len(files) == 1 {
			file := files[0]
			reader, err := file.Open()
			if err != nil {
				slog.Error("failed to open file", sl.Err(err))
				return c.JSON(echo.ErrInternalServerError.Code, throw("internal server error"))
			}
			defer reader.Close()

			e := domain.NewEntry(file.Filename, reader)
			path, err := svc.Upload(ctx, e)
			if err != nil {
				slog.Error("cannot upload file", sl.Err(err))
				return c.JSON(echo.ErrBadRequest.Code, throw("cannot upload file"))
			}
			paths = append(paths, path)
		} else {
			ee := make([]*domain.Entry, 0, len(files))
			for _, f := range files {
				filename := f.Filename
				reader, err := f.Open()
				if err != nil {
					slog.Error("failed to open file", sl.Err(err))
					return c.JSON(echo.ErrInternalServerError.Code, throw("internal server error"))
				}

				e := domain.NewEntry(filename, reader)
				ee = append(ee, e)
			}
			pp, err := svc.UploadMany(ctx, ee)
			if err != nil {
				slog.Error("cannot upload file", sl.Err(err))
				return c.JSON(echo.ErrBadRequest.Code, throw("cannot upload file"))
			}
			paths = append(paths, pp...)
		}

		return c.JSON(200, &H{
			"paths": paths,
		})
	}
}

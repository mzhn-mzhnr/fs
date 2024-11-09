package handlers

import (
	"bytes"
	"io"
	"log/slog"
	"mzhn/fileservice/internal/services/fileservice"
	"mzhn/fileservice/pkg/sl"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func File(svc *fileservice.FileService) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		id := c.QueryParam("id")
		if id == "" {
			return c.JSON(echo.ErrBadRequest.Code, throw("id is required"))
		}

		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(echo.ErrBadRequest.Code, throw("id must be uuid"))
		}

		file, err := svc.File(ctx, id)
		if err != nil {
			slog.Error("failed to get file", slog.String("id", id), sl.Err(err))
			return c.JSON(echo.ErrInternalServerError.Code, throw("failed to get file"))
		}

		body := new(bytes.Buffer)
		if _, err := io.Copy(body, file); err != nil {
			slog.Error("failed to read file", sl.Err(err))
			return c.JSON(echo.ErrInternalServerError.Code, throw("failed to read file"))
		}

		ct := http.DetectContentType(body.Bytes())

		return c.Blob(200, ct, body.Bytes())
	}
}

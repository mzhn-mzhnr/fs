package http

import (
	"context"
	"fmt"
	"log/slog"
	"mzhn/fileservice/internal/config"
	"mzhn/fileservice/internal/services/fileservice"
	"mzhn/fileservice/internal/transport/http/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Adapter struct {
	*echo.Echo
	svc *fileservice.FileService
	cfg *config.Config
}

func New(cfg *config.Config, svc *fileservice.FileService) *Adapter {
	return &Adapter{
		Echo: echo.New(),
		svc:  svc,
		cfg:  cfg,
	}
}

func (adapter *Adapter) setup() {
	e := adapter.Echo
	e.Use(middleware.Logger())
	e.Use(middleware.BodyLimit(adapter.cfg.Http.BodyLimit))

	port := adapter.cfg.Http.Port

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{fmt.Sprintf("http://localhost:%d", port)},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
	}))

	e.POST("/upload", handlers.Upload(adapter.svc))
	e.GET("/file/:filename", handlers.File(adapter.svc))
}

func (adapter *Adapter) Run(ctx context.Context) error {
	adapter.setup()

	host := adapter.cfg.Http.Host
	port := adapter.cfg.Http.Port
	addr := fmt.Sprintf("%s:%d", host, port)
	slog.Info("running http server", slog.String("addr", addr))

	go func() {
		if err := adapter.Start(addr); err != nil {
			return
		}
	}()

	<-ctx.Done()
	if err := adapter.Shutdown(ctx); err != nil {
		return err
	}

	slog.Info("shutting down http server\n")
	return nil
}

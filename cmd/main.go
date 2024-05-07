package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Aldikon/test_task_family_team/internal/controller"
	"github.com/Aldikon/test_task_family_team/internal/repository"
	"github.com/Aldikon/test_task_family_team/internal/service"
	"github.com/Aldikon/test_task_family_team/pkg/unzip"
	"github.com/go-chi/chi/v5"
)

func init() {
	flag.StringVar(&filePath, "file", "", "path to archive")
	flag.StringVar(&ext, "ext", "", "file extension")
	flag.StringVar(&port, "port", "8080", "service start port, default is 8080")

	flag.Parse()

	if filePath == "" {
		log.Fatal("path to the archive is mandatory, use the --file flag to specify the path")
	}
}

var (
	ext      string
	filePath string
	port     string
)

const (
	destDir       = "temp"
	timeOutServer = 2
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	slog.Info("path to archive", "path", filePath)
	slog.Info("file extension for output", "ext", ext)

	if err := os.Mkdir(destDir, os.ModeDir|0755); err != nil {

		if errors.Is(err, os.ErrExist) {
			if err := os.RemoveAll(filepath.Join(destDir, "/*")); err != nil {
				slog.Error("remove in dir", "path", destDir, "err", err)
				return
			}
		} else {
			slog.Error("create dir", "path", destDir, "err", err)
			return
		}
	}
	defer func() {
		if err := os.RemoveAll(destDir); err != nil {
			slog.Error("remove temp dir", "path", destDir, "err", err)
			return
		}
	}()

	if err := unzip.UnzipFile(filePath, destDir); err != nil {
		slog.Error("unzip file", "err", err)
		return
	}

	fsDir, err := os.ReadDir(destDir)
	if err != nil {
		slog.Error("read dir, err: %w", err)
		return
	}

	fsDirNew := make([]fs.DirEntry, 0, len(fsDir))

	for _, f := range fsDir {
		if ext == "" || filepath.Ext(f.Name()) == ext {
			fsDirNew = append(fsDirNew, f)
		}
	}

	r := chi.NewRouter()

	filerepo := repository.NewFile(fsDirNew, destDir)

	fileService := service.NewFile(filerepo)

	controller.New(r, ext, fileService)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	go func() {
		slog.Info("start work server", "adres", fmt.Sprintf("localhost:%s", port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Warn("run server", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), timeOutServer*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Warn("server shutdown", "err", err)
		return
	}

	<-ctx.Done()
	slog.Info("timeout of 2 seconds.")
	slog.Info("server exiting")
}

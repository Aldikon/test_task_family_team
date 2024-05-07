package controller

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type FileService interface {
	Get(ctx context.Context) []string
	GetByName(ctx context.Context, name string) (*os.File, error)
}

type file struct {
	ext         string
	fileService FileService
}

func New(r *chi.Mux, ext string, fileService FileService) {
	h := file{
		fileService: fileService,
		ext:         ext,
	}

	r.Get("/", h.Get)
	r.Get("/{name}", h.GetByName)
}

func (h *file) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res := h.fileService.Get(ctx)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Request{
		"ext":   h.ext,
		"names": res,
	})
}

func (h *file) GetByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	f, err := h.fileService.GetByName(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", http.DetectContentType([]byte{}))

	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

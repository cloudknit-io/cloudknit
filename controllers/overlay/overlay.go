package overlay

import (
	"github.com/compuzest/zlifecycle-il-operator/controllers/filereconciler"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
)

type Handler struct {
	fileService file.Service
	path        string
}

func NewHandler(fs file.Service, filepath string) *Handler {
	return &Handler{fileService: fs, path: filepath}
}

func (h *Handler) Reconcile() error {
	return nil
}

func (h *Handler) Cleanup() error {
	return h.fileService.RemoveAll(h.path)
}

var _ filereconciler.Handler = &Handler{}

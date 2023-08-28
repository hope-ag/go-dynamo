package health

import (
	"errors"
	"net/http"

	"github.com/hope-ag/go-dynamo/core/handlers"
	"github.com/hope-ag/go-dynamo/core/repository/adapter"
	"github.com/hope-ag/go-dynamo/utils/request"
)

type Handler struct {
	handlers.Interface
	Repository adapter.Interface
}

func NewHandler(repository adapter.Interface) handlers.Interface {
	return &Handler{
		Repository: repository,
	}
}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if !h.Repository.Health() {
		request.SendErrorResponse(w, r, errors.New("database error"), http.StatusInternalServerError)
		return
	}
	request.SendSuccessResponse(w, r, "All good")
}
func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	request.SendErrorResponse(w, r, nil, http.StatusMethodNotAllowed)
}
func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	request.SendErrorResponse(w, r, nil, http.StatusMethodNotAllowed)
}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	request.SendErrorResponse(w, r, nil, http.StatusMethodNotAllowed)
}
func (h *Handler) Patch(w http.ResponseWriter, r *http.Request) {
	request.SendErrorResponse(w, r, nil, http.StatusMethodNotAllowed)
}
func (h *Handler) Options(w http.ResponseWriter, r *http.Request) {
	request.SendErrorResponse(w, r, nil, http.StatusMethodNotAllowed)
}

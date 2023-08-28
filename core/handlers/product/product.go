package health

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/hope-ag/go-dynamo/core/controllers/product"
	EntityProduct "github.com/hope-ag/go-dynamo/core/entities/product"
	"github.com/hope-ag/go-dynamo/core/handlers"
	"github.com/hope-ag/go-dynamo/core/repository/adapter"
	"github.com/hope-ag/go-dynamo/core/rules"
	RulesProduct "github.com/hope-ag/go-dynamo/core/rules/product"
	"github.com/hope-ag/go-dynamo/utils/request"
)

type Handler struct {
	handlers.Interface
	Controller product.Interface
	Rules      rules.Interface
}

func NewHandler(repository adapter.Interface) handlers.Interface {
	return &Handler{
		Controller: product.NewController(repository),
		Rules:      RulesProduct.NewRules(),
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if chi.URLParam(r, "ID") != "" {
		h.getOne(w, r)
	} else {
		h.getAll(w, r)
	}
}

func (h *Handler) getOne(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		request.SendErrorResponse(w, r, errors.New("invalid ID"), 401)
		return
	}
	res, err := h.Controller.ListOne(ID)
	if err != nil {
		request.SendInternalError(w, r, err)
		return
	}
	request.SendSuccessResponse(w, r, res)
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	res, err := h.Controller.ListAll()
	if err != nil {
		request.SendInternalError(w, r, err)
		return
	}
	request.SendSuccessResponse(w, r, res)
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		request.SendBadRequest(w, r, errors.New("invalid id"))
		return
	}
	productBody, err := h.getBodyAndValidate(r, ID)

	if err != nil {
		request.SendBadRequest(w, r, err)
		return
	}

	if err := h.Controller.Update(ID, productBody); err != nil {
		request.SendInternalError(w, r, err)
		return
	}
	request.SendNoContentSuccessResponse(w, r)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	productBody, err := h.getBodyAndValidate(r, uuid.Nil)
	if err != nil {
		request.SendBadRequest(w, r, err)
		return
	}
	ID, err := h.Controller.Create(productBody)
	if err != nil {
		request.SendErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	request.SendSuccessResponse(w, r, map[string]any{"id": ID.String()})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		request.SendBadRequest(w, r, errors.New("invalid ID"))
		return
	}
	if err := h.Controller.Delete(ID); err != nil {
		request.SendInternalError(w, r, err)
		return
	}
	request.SendNoContentSuccessResponse(w, r)
}

func (h *Handler) Patch(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) Options(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) getBodyAndValidate(r *http.Request, ID uuid.UUID) (*EntityProduct.Product, error) {
	productBody := &EntityProduct.Product{}
	body, err := h.Rules.ConvertIoReaderToStruct(r.Body, productBody)
	if err != nil {
		return &EntityProduct.Product{}, err
	}

	productParsed, err := EntityProduct.InterfaceToModel(body)
	if err != nil {
		return &EntityProduct.Product{}, err
	}
	setDefaultValues(productParsed, ID)
	return productParsed, err
}

func setDefaultValues(product *EntityProduct.Product, ID uuid.UUID) {
	if ID == uuid.Nil {
		product.ID = uuid.New()
		product.CreatedAt = time.Now()
	} else {
		product.ID = ID
		product.UpdatedAt = time.Now()
	}
}

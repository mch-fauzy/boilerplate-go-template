package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/evermos/boilerplate-go/internal/domain/foobarbaz"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

// FooBarBazHandler is the HTTP handler for FooBarBaz domain.
type FooBarBazHandler struct {
	FooService     foobarbaz.FooService
	AuthMiddleware *middleware.Authentication
}

// ProvideFooBarBazHandler is the provider for this handler.
func ProvideFooBarBazHandler(fooService foobarbaz.FooService, authMiddleware *middleware.Authentication) FooBarBazHandler {
	return FooBarBazHandler{
		FooService:     fooService,
		AuthMiddleware: authMiddleware,
	}
}

// Router sets up the router for this domain.
func (h *FooBarBazHandler) Router(r chi.Router) {
	r.Route("/foobarbaz", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.AuthMiddleware.ClientCredential)
			r.Get("/foo/{id}", h.ResolveFooByID)
		})

		r.Group(func(r chi.Router) {
			r.Use(h.AuthMiddleware.Password)
			r.Post("/foo", h.CreateFoo)
			r.Delete("/foo/{id}", h.SoftDeleteFoo)
			r.Put("/foo/{id}", h.UpdateFoo)
		})

	})
}

// CreateFoo creates a new Foo.
// @Summary Create a new Foo.
// @Description This endpoint creates a new Foo.
// @Tags foobarbaz/foo
// @Security EVMOauthToken
// @Param foo body foobarbaz.FooRequestFormat true "The Foo to be created."
// @Produce json
// @Success 201 {object} response.Base{data=foobarbaz.FooResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/foobarbaz/foo [post]
func (h *FooBarBazHandler) CreateFoo(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat foobarbaz.FooRequestFormat
	err := decoder.Decode(&requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID, _ := uuid.NewV4() // TODO: read from context

	foo, err := h.FooService.Create(requestFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, foo)
}

// ResolveFooByID resolves a Foo by its ID.
// @Summary Resolve Foo by ID
// @Description This endpoint resolves a Foo by its ID.
// @Tags foobarbaz/foo
// @Security EVMOauthToken
// @Param id path string true "The Foo's identifier."
// @Param withItems query string false "Fetch with items, default false."
// @Produce json
// @Success 200 {object} response.Base{data=foobarbaz.FooResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 404 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/foobarbaz/foo/{id} [get]
func (h *FooBarBazHandler) ResolveFooByID(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, err := uuid.FromString(idString)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	withItems, _ := strconv.ParseBool(r.URL.Query().Get("withItems"))

	foo, err := h.FooService.ResolveByID(id, withItems)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, foo)
}

// SoftDeleteFoo marks a Foo as deleted.
// @Summary Marks a Foo as deleted.
// @Description This endpoint marks an existing Foo as deleted. This is done by
// @Description setting the "deleted" and "deletedBy" properties of the Foo.
// @Tags foobarbaz/foo
// @Security EVMOauthToken
// @Param id path string true "The Foo's identifier."
// @Produce json
// @Success 200 {object} response.Base{data=foobarbaz.FooResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/foobarbaz/foo/{id} [delete]
func (h *FooBarBazHandler) SoftDeleteFoo(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, err := uuid.FromString(idString)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID, _ := uuid.NewV4() // TODO: read from context

	foo, err := h.FooService.SoftDelete(id, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, foo)
}

// UpdateFoo updates a Foo.
// @Summary Update a Foo.
// @Description This endpoint updates an existing Foo.
// @Tags foobarbaz/foo
// @Security EVMOauthToken
// @Param id path string true "The Foo's identifier."
// @Param foo body foobarbaz.FooRequestFormat true "The Foo to be updated."
// @Produce json
// @Success 200 {object} response.Base{data=foobarbaz.FooResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/foobarbaz/foo/{id} [put]
func (h *FooBarBazHandler) UpdateFoo(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, err := uuid.FromString(idString)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var requestFormat foobarbaz.FooRequestFormat
	err = decoder.Decode(&requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID, _ := uuid.NewV4() // TODO: read from context

	foo, err := h.FooService.Update(id, requestFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, foo)
}

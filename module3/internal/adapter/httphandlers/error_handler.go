package httphandlers

import (
	"errors"
	"module3/internal/repository"
	"net/http"

	"github.com/go-chi/render"
)

func RespondError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, repository.ErrAliasAlreadyExists):
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, map[string]string{"error": "Alias already exists"})

	case errors.Is(err, repository.ErrNotFound):
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "Not found"})

	case errors.Is(err, repository.ErrUnauthorized):
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "Unauthorized"})

	default:
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Internal server error"})
	}
}

func RenderInvalidJsonError(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, map[string]string{"error": "Invalid JSON format"})
}

func RenderInvalidIdError(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, map[string]string{"error": "Invalid ID format"})
}

func RenderBadRequestError(w http.ResponseWriter, r *http.Request, errorMessage string) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, map[string]string{"error": errorMessage})
}

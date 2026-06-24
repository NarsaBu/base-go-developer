package httphandlers

import (
	"encoding/json"
	"module3/internal/dto"
	"module3/internal/usecase"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type UrlHandler struct {
	urlService *usecase.UrlService
}

func NewUrlHandler(service *usecase.UrlService) *UrlHandler {
	return &UrlHandler{urlService: service}
}

func (uh *UrlHandler) HandleCreateUrl(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Url   string `json:"url"`
		Alias string `json:"alias"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RenderInvalidJsonError(w, r)
		return
	}

	if req.Url == "" || req.Alias == "" {
		RenderBadRequestError(w, r, "Fields 'url' and 'alias' are required and cannot be empty")
		return
	}

	createdUrl, err := uh.urlService.CreateUrl(req.Url, req.Alias)
	if err != nil {
		RespondError(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, createdUrl)
}

func (uh *UrlHandler) HandleUpdateUrl(w http.ResponseWriter, r *http.Request) {
	var req dto.UrlUpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RenderInvalidJsonError(w, r)
		return
	}

	if req.Id == 0 || req.Url == "" || req.Alias == "" {
		RenderBadRequestError(w, r, "Fields 'id', 'url' and 'alias' are required and cannot be empty")
		return
	}

	updatedUrl, err := uh.urlService.UpdateUrl(&req)
	if err != nil {
		RespondError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, updatedUrl)
}

func (uh *UrlHandler) HandleDeleteUrlById(w http.ResponseWriter, r *http.Request) {
	urlId := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(urlId, 10, 64)
	if err != nil {
		RenderInvalidIdError(w, r)
		return
	}

	err = uh.urlService.DeleteById(id)
	if err != nil {
		RespondError(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

func (uh *UrlHandler) HandleGetUrlById(w http.ResponseWriter, r *http.Request) {
	urlId := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(urlId, 10, 64)
	if err != nil {
		RenderInvalidIdError(w, r)
		return
	}

	foundUrl, err := uh.urlService.FindById(id)
	if err != nil {
		RespondError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, foundUrl)
}

func (uh *UrlHandler) HandleRedirectByAlias(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")
	if alias == "" {
		RenderBadRequestError(w, r, "Alias should not be empty")
		return
	}

	urlToRedirect, err := uh.urlService.FindUrlStringByAlias(alias)
	if err != nil {
		RespondError(w, r, err)
		return
	}

	http.Redirect(w, r, urlToRedirect, http.StatusFound)
}

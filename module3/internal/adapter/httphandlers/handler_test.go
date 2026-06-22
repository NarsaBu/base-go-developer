package httphandlers

import (
	"bytes"
	"io"
	"module3/internal/adapter/postgres/mocks"
	"module3/internal/entities"
	"module3/internal/repository"
	"module3/internal/usecase"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	_ = godotenv.Load("../../.env")
}

func setupTestRouter(mockRepo *mocks.MockUrlRepository) *chi.Mux {
	authUser := os.Getenv("AUTH_USERNAME")
	authPass := os.Getenv("AUTH_PASSWORD")

	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok || user != authUser || pass != authPass {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	urlService := usecase.NewUrlService(mockRepo)
	urlHandler := NewUrlHandler(urlService)

	r.Route("/urls", func(r chi.Router) {
		r.Post("/", urlHandler.HandleCreateUrl)
		r.Put("/", urlHandler.HandleUpdateUrl)
		r.Delete("/{id}", urlHandler.HandleDeleteUrlById)
		r.Get("/{id}", urlHandler.HandleGetUrlById)
		r.Get("/redirect/{alias}", urlHandler.HandleRedirectByAlias)
	})

	return r
}

func makeRequest(r *chi.Mux, method, path string, body io.Reader) *http.Response {
	req := httptest.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")

	authUser := os.Getenv("AUTH_USERNAME")
	authPass := os.Getenv("AUTH_PASSWORD")
	req.SetBasicAuth(authUser, authPass)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr.Result()
}

func TestCreateUrl_Success(t *testing.T) {
	mockRepo := new(mocks.MockUrlRepository)
	router := setupTestRouter(mockRepo)

	reqBody := `{"url": "https://example.com", "alias": "ex"}`
	mockRepo.On("Save", "https://example.com", "ex").Return(&entities.Url{
		Id: 1, Url: "https://example.com", Alias: "ex",
	}, nil)

	resp := makeRequest(router, http.MethodPost, "/urls/", bytes.NewBufferString(reqBody))

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestCreateUrl_BadRequest(t *testing.T) {
	mockRepo := new(mocks.MockUrlRepository)
	router := setupTestRouter(mockRepo)

	reqBody := `{invalid json!!!}`

	resp := makeRequest(router, http.MethodPost, "/urls/", bytes.NewBufferString(reqBody))

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mockRepo.AssertNotCalled(t, "Save")
}

func TestCreateUrl_Conflict(t *testing.T) {
	mockRepo := new(mocks.MockUrlRepository)
	router := setupTestRouter(mockRepo)

	reqBody := `{"url": "https://example.com", "alias": "ex"}`
	mockRepo.On("Save", "https://example.com", "ex").Return(nil, repository.ErrAliasAlreadyExists)

	resp := makeRequest(router, http.MethodPost, "/urls/", bytes.NewBufferString(reqBody))

	assert.Equal(t, http.StatusConflict, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUrl_Success(t *testing.T) {
	mockRepo := new(mocks.MockUrlRepository)
	router := setupTestRouter(mockRepo)

	reqBody := `{"id": 1, "url": "https://updated.com", "alias": "upd"}`
	mockRepo.On("Update", mock.MatchedBy(func(u *entities.Url) bool {
		return u.Id == 1 && u.Url == "https://updated.com" && u.Alias == "upd"
	})).Return(&entities.Url{
		Id: 1, Url: "https://updated.com", Alias: "upd",
	}, nil)

	resp := makeRequest(router, http.MethodPut, "/urls/", bytes.NewBufferString(reqBody))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUrl_BadRequest(t *testing.T) {
	mockRepo := new(mocks.MockUrlRepository)
	router := setupTestRouter(mockRepo)

	reqBody := `{broken json!!!}`
	resp := makeRequest(router, http.MethodPut, "/urls/", bytes.NewBufferString(reqBody))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdateUrl_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockUrlRepository)
	router := setupTestRouter(mockRepo)

	reqBody := `{"id": 999, "url": "https://updated.com", "alias": "upd"}`
	mockRepo.On("Update", mock.Anything).Return(nil, repository.ErrNotFound)

	resp := makeRequest(router, http.MethodPut, "/urls/", bytes.NewBufferString(reqBody))
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUrl_Success(t *testing.T) {
	mockRepo := new(mocks.MockUrlRepository)
	router := setupTestRouter(mockRepo)

	mockRepo.On("DeleteById", int64(1)).Return(nil)

	resp := makeRequest(router, http.MethodDelete, "/urls/1", nil)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUrl_InvalidId(t *testing.T) {
	mockRepo := new(mocks.MockUrlRepository)
	router := setupTestRouter(mockRepo)

	resp := makeRequest(router, http.MethodDelete, "/urls/abc", nil)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mockRepo.AssertNotCalled(t, "DeleteById")
}

func TestGetUrlById_Success(t *testing.T) {
	mockRepo := new(mocks.MockUrlRepository)
	router := setupTestRouter(mockRepo)

	mockRepo.On("FindById", int64(1)).Return(&entities.Url{
		Id: 1, Url: "https://example.com", Alias: "ex",
	}, nil)

	resp := makeRequest(router, http.MethodGet, "/urls/1", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestRedirectByAlias_Success(t *testing.T) {
	mockRepo := new(mocks.MockUrlRepository)
	router := setupTestRouter(mockRepo)

	mockRepo.On("FindUrlStringByAlias", "ex").Return("https://example.com", nil)

	req := httptest.NewRequest(http.MethodGet, "/urls/redirect/ex", nil)
	authUser := os.Getenv("AUTH_USERNAME")
	authPass := os.Getenv("AUTH_PASSWORD")
	req.SetBasicAuth(authUser, authPass)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusFound, rr.Code)
	assert.Equal(t, "https://example.com", rr.Header().Get("Location"))
	mockRepo.AssertExpectations(t)
}

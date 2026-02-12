package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/", h.Shorten)
	r.GET("/:id", h.Redirect)
	return r
}

func TestGenerateRandomString_Length(t *testing.T) {
	h := NewHandler()
	id := h.generateRandomString()
	assert.Len(t, id, 6)
}

func TestShorten_Success(t *testing.T) {
	h := NewHandler()
	router := setupRouter(h)

	body := "https://practicum.yandex.ru/"
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "http://localhost:8080/")
	assert.Len(t, h.urlMap, 1)
}

func TestShorten_InvalidMethod(t *testing.T) {
	h := NewHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestShorten_InvalidURL(t *testing.T) {
	h := NewHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("invalid-url"))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRedirect_Success(t *testing.T) {
	h := NewHandler()
	h.urlMap["genstr"] = "https://practicum.yandex.ru/"

	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/genstr", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, "https://practicum.yandex.ru/", w.Header().Get("Location"))
}

func TestRedirect_UnknownID(t *testing.T) {
	h := NewHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRedirect_InvalidMethod(t *testing.T) {
	h := NewHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/example", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

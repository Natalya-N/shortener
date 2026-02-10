package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomString_Length(t *testing.T) {
	id := generateRandomString()
	assert.Len(t, id, 6)
}

func TestShortenHandler_Success(t *testing.T) {

	urlMap = make(map[string]string)

	body := "https://practicum.yandex.ru/"
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	rec := httptest.NewRecorder()

	shortenHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "text/plain", resp.Header.Get("Content-Type"))

	responseBody := rec.Body.String()
	assert.Contains(t, responseBody, "http://localhost:8080/")
	assert.Len(t, urlMap, 1)
}

func TestShortenHandler_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	shortenHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestShortenHandler_InvalidURL(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("invalid-url"))
	rec := httptest.NewRecorder()

	shortenHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRedirectHandler_Success(t *testing.T) {
	urlMap = make(map[string]string)
	urlMap["genstr"] = "https://practicum.yandex.ru/"

	req := httptest.NewRequest(http.MethodGet, "/genstr", nil)
	rec := httptest.NewRecorder()

	redirectHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.Equal(t, "https://practicum.yandex.ru/", resp.Header.Get("Location"))
}

func TestRedirectHandler_UnknownID(t *testing.T) {
	urlMap = make(map[string]string)

	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	rec := httptest.NewRecorder()

	redirectHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRedirectHandler_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/example", nil)
	rec := httptest.NewRecorder()

	redirectHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

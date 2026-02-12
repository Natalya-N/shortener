package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	urlMap map[string]string
	rng    *rand.Rand
}

func NewHandler() *Handler {
	return &Handler{
		urlMap: make(map[string]string),
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (h *Handler) generateRandomString() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var b strings.Builder
	for i := 0; i < 6; i++ {
		b.WriteByte(charset[h.rng.Intn(len(charset))])
	}
	return b.String()
}

func (h *Handler) Shorten(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil || len(body) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	originalURL := strings.TrimSpace(string(body))
	parsedURL, err := url.Parse(originalURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	id := h.generateRandomString()
	h.urlMap[id] = originalURL

	shortURL := fmt.Sprintf("http://localhost:8080/%s", id)

	c.Header("Content-Type", "text/plain")
	c.String(http.StatusCreated, shortURL)
}

func (h *Handler) Redirect(c *gin.Context) {
	id := c.Param("id")

	if id == "" || strings.Contains(id, "/") {
		c.Status(http.StatusBadRequest)
		return
	}

	originalURL, ok := h.urlMap[id]
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	c.Header("Location", originalURL)
	c.Status(http.StatusTemporaryRedirect)
}

package main

import (
	"fmt"

	"github.com/Natalya-N/shortener/internal/config"
	"github.com/Natalya-N/shortener/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.NewConfig()

	fmt.Printf("Server started on %s\n", cfg.ServerAddress)

	h := handler.NewHandler(cfg.BaseURL)

	r := gin.Default()

	r.POST("/", h.Shorten)
	r.GET("/:id", h.Redirect)

	r.Run(cfg.ServerAddress)
}

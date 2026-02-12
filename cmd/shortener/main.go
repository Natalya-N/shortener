package main

import (
	"fmt"

	"github.com/Natalya-N/shortener/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Server started on http://localhost:8080")

	h := handler.NewHandler()

	r := gin.Default()

	r.POST("/", h.Shorten)
	r.GET("/:id", h.Redirect)

	r.Run(":8080")
}

package main

import (
	"log"
	"net/http"
	"os"

	"les-randoms/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.Foo(1)

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})
	router.GET("/lol", func(c *gin.Context) {
		c.HTML(http.StatusOK, "lol-index.tmpl.html", nil)
	})

	router.Run(":" + port)
}

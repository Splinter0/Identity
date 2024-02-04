package main

import (
	"github.com/Splinter0/identity/endpoints"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/js/", "static/js/")
	config := endpoints.LoadConfig()
	for _, provider := range config.Providers {
		if provider == "bankid" {
			endpoints.RegisterBankIDEndpoints(r, config)
		}
	}
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"Providers": config.Providers,
		})
	})
	r.Run("0.0.0.0:8080")
}

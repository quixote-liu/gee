package main

import (
	"net/http"

	"gee"
)

func main() {
	e := gee.New()
	e.GET("/get", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"method": c.Method,
			"path":   c.Path,
			"value":  "hello, world",
		})
	})
	e.POST("/post", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"method": c.Method,
			"path":   c.Path,
			"value":  "hello, world",
		})
	})
	e.PUT("/put", func(c *gee.Context) {
		c.String(http.StatusOK, "hello, put method, path: %s", c.Method)
	})
	e.Run(":8080")
}

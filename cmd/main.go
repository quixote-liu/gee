package main

import (
	"net/http"

	"gee"
)

func main() {
	e := gee.New()
	v1 := e.Group("/v1")
	{
		v1.GET("/user", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"hello": "lcs",
			})
		})
		v1.POST("/client", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"hello": "client",
			})
		})
	}

	v2 := e.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			name := c.Param("name")
			c.JSON(http.StatusOK, gee.H{
				"hello": name,
			})
		})
		v2.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"hello": "v2 post",
			})
		})
	}

	e.GET("/v3/hello/:name", func(c *gee.Context) {
		name := c.Param("name")
		c.JSON(http.StatusOK, gee.H{
			"hello": name,
		})
	})

	e.Run(":8080")
}

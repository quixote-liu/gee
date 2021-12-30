package main

import (
	"encoding/xml"
	"net/http"

	"gee"
)

func main() {
	e := gee.New()

	e.Use(gee.Logger())
	e.LoadHTMLGlob("./templates/*")

	delivery := gee.H{
		"name":         "lcs",
		"id":           "2020",
		"address":      "hubeiwuhan",
		"organization": "hubeidaxue",
		"hobby":        "write code",
		"profession":   "programmer",
	}

	// JSON
	e.GET("/hello", func(c *gee.Context) {
		c.JSON(http.StatusOK, delivery)
	})

	// secureJSON
	e.SecureJsonPrefix("hello")
	e.GET("/securejson", func(c *gee.Context) {
		c.SecureJSON(http.StatusOK, delivery)
	})

	// IndentedJSON
	e.GET("/indentedjson", func(c *gee.Context) {
		c.IndentedJSON(http.StatusOK, delivery)
	})

	// JSONP
	e.GET("/jsonp", func(c *gee.Context) {
		c.JSONP(http.StatusOK, delivery)
	})

	// AscillJSON
	e.GET("/asciijson", func(c *gee.Context) {
		c.AsciiJSON(http.StatusOK, delivery)
	})

	// XML
	e.GET("/xml", func(c *gee.Context) {
		type server struct {
			ServerName string `xml:"serverName"`
			ServerIP   string `xml:"serverIP"`
		}
		type Servers struct {
			XMLName xml.Name `xml:"servers"`
			Version string   `xml:"version,attr"`
			Svs     []server `xml:"server"`
		}
		v := &Servers{Version: "1"}
		v.Svs = append(v.Svs, server{"Shanghai_VPN", "127.0.0.1"})
		v.Svs = append(v.Svs, server{"Beijing_VPN", "127.0.0.2"})

		c.XML(http.StatusOK, v)
	})

	// YAML
	e.GET("/yaml", func(c *gee.Context) {
		c.YAML(http.StatusOK, delivery)
	})

	// String
	e.GET("/string", func(c *gee.Context) {
		c.String(http.StatusOK, "hello, lcs %s", "heiheihei")
	})

	// HTML
	e.GET("/html", func(c *gee.Context) {
		type todo struct {
			Title string
			Done  bool
		}
		type todoPageData struct {
			PageTitle string
			Todos     []todo
		}
		c.HTML(http.StatusOK, "temp1", todoPageData{
			PageTitle: "My TODO list",
			Todos: []todo{
				{Title: "Task 1", Done: false},
				{Title: "Task 2", Done: true},
				{Title: "Task 3", Done: true},
			},
		})
	})

	// e.GET("/", func(c *gee.Context) {
	// 	c.HTML(http.StatusOK, "css.tmpl", nil)
	// })

	// e.Static("/assert", "./static/")

	// v1 := e.Group("/v1")
	// {
	// 	v1.GET("/user", func(c *gee.Context) {
	// 		c.JSON(http.StatusOK, gee.H{
	// 			"hello": "lcs",
	// 		})
	// 	})
	// 	v1.POST("/client", func(c *gee.Context) {
	// 		c.JSON(http.StatusOK, gee.H{
	// 			"hello": "client",
	// 		})
	// 	})
	// }

	// v2 := e.Group("/v2")
	// {
	// 	v2.GET("/hello/:name", func(c *gee.Context) {
	// 		name := c.Param("name")
	// 		c.JSON(http.StatusOK, gee.H{
	// 			"hello": name,
	// 		})
	// 	})
	// 	v2.POST("/login", func(c *gee.Context) {
	// 		c.JSON(http.StatusOK, gee.H{
	// 			"hello": "v2 post",
	// 		})
	// 	})
	// }

	// e.GET("/v3/hello/:name", func(c *gee.Context) {
	// 	name := c.Param("name")
	// 	c.JSON(http.StatusOK, gee.H{
	// 		"hello": name,
	// 	})
	// })

	e.Run(":8080")
}

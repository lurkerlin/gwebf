package main

import (
	"fmt"
	"html/template"
	"log"
	"lweb"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age  int
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := lweb.New()
	r.Use(lweb.Logger()) // global middleware
	r.Use(lweb.Recovery())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	v1 := r.Group("/v1")
	{
		v1.GET("/hello", func(c *lweb.Context) {
			c.String(http.StatusOK, "hello v1 group")
		})
	}
	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) //group middleware
	{
		v2.GET("/hello", func(c *lweb.Context) {
			c.String(http.StatusOK, "v2 group")
		})
	}
	r.GET("/", func(c *lweb.Context) {
		c.HTML(http.StatusOK, "lweb.html", nil)
	})
	r.GET("/hello", func(c *lweb.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *lweb.Context) {
		c.JSON(http.StatusOK, lweb.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.GET("/article/:id", func(c *lweb.Context) {
		c.String(http.StatusOK, "article id: %s\n", c.Param("id"))
	})
	r.GET("/panic", func(c *lweb.Context) {
		names := []string{"hello world"}
		c.String(http.StatusOK, names[100])
	})
	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.GET("/students", func(c *lweb.Context) {
		c.HTML(http.StatusOK, "arr.html", lweb.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})
	r.GET("/date", func(c *lweb.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", lweb.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})
	r.Static("/assets", "./static")
	r.Run(":8080")
}

func onlyForV2() lweb.HandlerFunc {
	return func(c *lweb.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

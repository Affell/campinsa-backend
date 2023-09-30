package handlers

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
)

func IndexHandler(c iris.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.StopWithJSON(iris.StatusOK, iris.Map{"message": "Welcome to the Oui API"})
}

func FaviconHandler(static embed.FS) func(c iris.Context) {
	return func(c iris.Context) {
		if c.RequestPath(false) != "/favicon.ico" {
			return
		}

		if method := c.Request().Method; method != "GET" && method != "HEAD" {
			status := iris.StatusOK
			if method != "OPTIONS" {
				status = iris.StatusMethodNotAllowed
			}
			c.Header("Allow", "GET,HEAD,OPTIONS")
			c.StopWithJSON(status, iris.Map{"error": status})
		}

		c.Header("Content-Type", "image/x-icon")

		favicon, err := static.ReadFile("public/favicon.ico")
		if err != nil {
			golog.Fatalf("init favicon error: %s", err)
		}

		c.Binary(favicon)
		c.StopWithStatus(iris.StatusOK)
	}
}

func FileHandler(static embed.FS, dir string) http.FileSystem {
	sub, err := fs.Sub(static, dir)

	if err != nil {
		golog.Fatalf("when accessing public path: %s | error : %s", dir, err)
	}

	return http.FS(sub)
}

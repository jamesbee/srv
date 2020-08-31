package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/russross/blackfriday/v2"
)

func (e *Engine) ServeFiles() *Engine {
	files := e.files
	if len(files) == 1 && len(e.dirs) == 0 {
		handler := e.doServeFile(files[0])

		e.srv.GET(genericURL(files[0]), handler)
		e.indexHandler = handler
	} else {
		for _, f := range files {
			fp := genericPath(f)
			if strings.HasSuffix(fp, "index") ||
				strings.HasSuffix(fp, "index.htm") ||
				strings.HasSuffix(fp, "index.html") {
				e.customIndex = true
			}
			e.srv.GET(genericURL(f), e.doServeFile(fp))
		}
	}
	return e
}

func (e *Engine) doServeFile(f string) echo.HandlerFunc {
	f = genericPath(f)
	if EnableMarkdown &&
		(strings.HasSuffix(f, "md") ||
			strings.HasSuffix(f, "markdown")) {
		return e.doServeMarkdown(f)
	}

	return func(c echo.Context) error {
		return c.File(f)
	}
}

func (e *Engine) doServeMarkdown(f string) echo.HandlerFunc {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}
	render := blackfriday.Run(data)
	return func(c echo.Context) error {
		return c.HTMLBlob(http.StatusOK, render)
	}
}

func (e *Engine) addFiles(files ...string) {
	e.files = append(e.files, files...)
}

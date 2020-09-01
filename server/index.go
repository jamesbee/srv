package server

import (
	"bytes"
	"html/template"
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
)

var html = template.Must(template.New("index").Parse(`
<html>
<head>
  <title>Server index</title>
</head>
<body>
{{ if .Index }}
	<h2>Index</h2>
	<ul>
	{{ range .Index }}
		<li><a href="{{ print . }}">{{ print . }}</a></li>
	{{ end }}
	</ul>
{{ end }}
{{ if .Files }}
	<h2>Files</h2>
	<ul>
	{{ range .Files }}
		<li><a href="/{{ print . }}">{{ print . }}</a></li>
	{{ end }}
	</ul>
{{ end }}
{{ if .Dirs }}
	<h2>Dirs</h2>
	<ul>
	{{ range .Dirs }}
		{{ if eq . "." }}
			<li><a href="/{{ print $.Static }}">/{{ print $.Static }}</a></li>
		{{ else }}
			<li><a href="/{{- print . }}">{{ print . }}</a></li>
		{{ end }}
	{{ end }}
	</ul>
{{ end }}
</body>
</html>
`))

func (e *Engine) setupIndex() {
	if e.indexHandler == nil {
		e.indexHandler = func(c echo.Context) (err error) {
			return indexRouteInfo(c, e.index, e.files, e.dirs)
		}
	}
	e.srv.GET("/", e.indexHandler)
	if !e.customIndex {
		e.srv.GET("/index", e.indexHandler)
		e.srv.GET("/index.html", e.indexHandler)
		e.index = []string{"/", "/index", "/index.html"}
	} else {
		e.index = []string{"/"}
	}
}

func indexRouteInfo(c echo.Context, index, files, dirs []string) (err error) {
	sort.Strings(index)
	sort.Strings(files)
	sort.Strings(dirs)
	buf := new(bytes.Buffer)
	err = html.ExecuteTemplate(buf, "index", H{
		"Index":  index,
		"Files":  files,
		"Dirs":   dirs,
		"Static": Static,
	})
	if err != nil {
		return
	}

	return c.HTMLBlob(http.StatusOK, buf.Bytes())
}

package server

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

func (e *Engine) ServeDirs() *Engine {
	dirs := e.dirs
	if len(dirs) == 1 {
		e.doServeDir(dirs[0])
	} else {
		for _, d := range dirs {
			e.doServeDir(d)
		}
	}
	return e
}

func (e *Engine) doServeDir(dp string) {
	srv := e.srv
	if dp[len(dp)-1] == '/' {
		dp = dp[:len(dp)-1]
	}
	dp = clarifyPath(dp)
	if dp == "." {
		srv.GET("/"+Static, dispatch())
		srv.GET("/"+Static+"/:uri", dispatch())
	} else {
		srv.GET("/"+dp, dispatch())
		srv.GET("/"+dp+"/:uri", dispatch())
	}
}

func dispatch() echo.HandlerFunc {
	return func(c echo.Context) error {
		uri := c.Param("uri")
		if uri == "" || uri == "/" {
			requestURI := strings.TrimPrefix(c.Request().RequestURI, "/"+Static)
			return listFile(c, genericPath(requestURI))
		}
		fs, err := os.Stat(uri)
		if err != nil {
			return err
		}
		return doDispatch(c, fs)
	}
}

func doDispatch(c echo.Context, fs os.FileInfo) (err error) {
	uri := genericPath(c.Request().RequestURI)
	if fs.IsDir() {
		return listFile(c, uri)
	}
	return c.File(uri)
}

func listFile(c echo.Context, uri string) (err error) {
	var files []string
	var dirs []string
	err = filepath.Walk(uri, func(path string, info os.FileInfo, err error) error {
		if path == "." || path == "./." || path == uri || isExclude(genericPath(path)) {
			return nil
		}
		if info.IsDir() {
			dirs = append(dirs, path)
		} else {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return
	}

	return indexRouteInfo(c, []string{}, files, dirs)
}

func (e *Engine) addDirs(dirs ...string) {
	e.dirs = append(e.dirs, dirs...)
}

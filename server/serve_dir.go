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
		srv.GET("/"+Static, e.dispatch(true))
		srv.GET("/"+Static+"/:uri", e.dispatch(true))
	} else {
		srv.GET("/"+dp, e.dispatch(false))
		srv.GET("/"+dp+"/:uri", e.dispatch(false))
	}
}

func (e *Engine) dispatch(bypass bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		// check if listing file action
		uri := c.Param("uri")
		if uri == "" || uri == "/" {
			requestURI := strings.TrimPrefix(c.Request().RequestURI, "/"+Static)
			return listFile(c, genericPath(requestURI))
		}

		var requestUri string
		if bypass {
			// serve current dir and bypass path prefix
			requestUri = c.Param("uri")
		} else {
			// serve file path
			requestUri = c.Request().RequestURI
		}

		fs, err := os.Stat(genericPath(requestUri))
		if err != nil {
			return err
		}
		return e.doDispatch(c, fs)
	}
}

func (e *Engine) doDispatch(c echo.Context, fs os.FileInfo) (err error) {
	uri := genericPath(c.Request().RequestURI)
	if fs.IsDir() {
		return listFile(c, uri)
	}
	return e.doServeFile(uri)(c)
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

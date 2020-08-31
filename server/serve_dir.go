package server

import (
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func (e *Engine) ServeDirs() *Engine {
	dirs := e.dirs
	if len(dirs) == 1 {
		doServeDir(e.srv, dirs[0])
	} else {
		for _, d := range dirs {
			doServeDir(e.srv, d)
		}
	}
	return e
}

func doServeDir(r *echo.Echo, dp string) {
	if dp[len(dp)-1] == '/' {
		dp = dp[:len(dp)-1]
	}
	if dp == "." {
		r.GET("/"+Static, dispatch())
		r.GET("/"+Static+"/:uri", dispatch())
	} else {
		r.GET("/"+dp, dispatch())
		r.GET("/"+dp+"/:uri", dispatch())
	}
}

func dispatch() echo.HandlerFunc {
	return func(c echo.Context) error {
		uri := c.Param("uri")
		if uri == "" || uri == "/" {
			return listFile(c, genericPath(c.Request().RequestURI))
		}
		return doDispatch(c)
	}
}

func doDispatch(c echo.Context) (err error) {
	uri := genericPath(c.Request().RequestURI)
	fs, err := os.Stat(uri)
	if err != nil {
		return
	}
	if fs.IsDir() {
		return listFile(c, uri)
	}
	return c.File(uri)
}

func listFile(c echo.Context, uri string) (err error) {
	var files []string
	var dirs []string
	err = filepath.Walk(uri, func(path string, info os.FileInfo, err error) error {
		if path == "." || path == uri {
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

	return sendFileList(c, []string{}, files, dirs)
}

func (e *Engine) addDirs(dirs ...string) {
	e.dirs = append(e.dirs, dirs...)
}

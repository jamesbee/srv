package server

import (
	"errors"
	"path/filepath"
	"strings"
)

type H map[string]interface{}

func genericPath(path string) string {
	f := filepath.Clean(path)
	if strings.HasPrefix(f, "..") {
		panic(errors.New("Path leak detected! should not access parent dir: " + path))
	}
	if f[0] == '/' {
		return "." + f
	} else if !strings.HasPrefix(f, "./") {
		return "./" + f
	}
	return f
}

func genericURL(path string) string {
	f := filepath.Clean(path)
	if strings.HasPrefix(f, "..") {
		panic(errors.New("Path leak detected! should not access parent dir: " + path))
	}
	if f[0] == '.' {
		f = f[1:]
	}
	if f[0] != '/' {
		f = "/" + f
	}
	if f[len(f)-1] == '/' {
		return f[:len(f)-1]
	}
	return f
}

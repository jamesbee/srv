package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

type Engine struct {
	index        []string
	files        []string
	dirs         []string
	srv          *echo.Echo
	indexHandler func(c echo.Context) (err error)
	customIndex  bool
}

func (e *Engine) Serve(source ...string) *Engine {
	if len(source) == 0 {
		e.addDirs(".")
		e.ServeDirs()
		return e
	}

	for _, f := range source {
		fi, err := os.Stat(f)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("File not found: ", f)
				os.Exit(0)
			}
			panic(err)
		}
		if fi.IsDir() {
			e.dirs = append(e.dirs, f)
		} else {
			e.files = append(e.files, f)
		}
	}
	if len(e.dirs) != 0 {
		e.ServeDirs()
	}
	if len(e.files) != 0 {
		e.ServeFiles()
	}

	return e
}

func (e *Engine) Listen(addr string) {
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := e.srv.Start(addr); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

func (e *Engine) ListenTLS(addr, cert, key string) {
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := e.srv.StartTLS(addr, cert, key); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

func (e *Engine) setup(addr string) {
	e.setupIndex()
	e.setupFallback()
}

func (e *Engine) setupFallback() {
	e.srv.HTTPErrorHandler = func(err error, c echo.Context) {
		if err != echo.ErrNotFound {
			return
		}
		uri := genericPath(c.Request().RequestURI)
		fs, err := os.Stat(uri)
		if err != nil || fs.IsDir() {
			c.Error(err)
			return
		}
		err = e.doServeFile(uri)(c)
		if err != nil {
			c.Error(err)
		}
	}
}

func (e *Engine) endless() {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

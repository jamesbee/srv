package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jamesbee/srv/config"
)

var (
	UseTLS         bool
	EnableMarkdown bool
	Port           int
	Host           string
	Static         string
	Cert           string
	Key            string
)

func CommandUp(cmd *cobra.Command, args []string) {
	e := NewEngine().
		Serve(args...)

	addr := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))
	e.setup(addr)

	if UseTLS {
		e.ListenTLS(addr, viper.GetString("tlsCert"), viper.GetString("tlsKey"))
	} else {
		e.Listen(addr)
	}

	e.endless()
}

func NewEngine() *Engine {
	srv := echo.New()
	srv.HideBanner = true
	if config.Debug {
		srv.Debug = true
	} else {
		srv.Debug = false
	}

	// remove trailing slash
	srv.Pre(middleware.Rewrite(map[string]string{
		"/*/": "/$1",
	}))
	// common plugin set
	srv.Use(middleware.Recover(),
		middleware.Logger(),
		middleware.Gzip())
	return &Engine{srv: srv}
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

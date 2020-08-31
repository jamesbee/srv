package server

import (
	"fmt"

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

	srv.Pre(middleware.Rewrite(map[string]string{
		"/*/": "/$1",
	}))
	srv.Use(middleware.Recover(),
		middleware.Logger(),
		middleware.Gzip())
	return &Engine{srv: srv}
}

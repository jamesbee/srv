package main

import (
	rice "github.com/GeertJohan/go.rice"

	"github.com/jamesbee/srv/cmd"
	"github.com/jamesbee/srv/config"
)

func main() {
	config.Assets = rice.MustFindBox("assets")
	cmd.Main()
}

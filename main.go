package main

import (
	"log"

	"github.com/jamesbee/srv/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

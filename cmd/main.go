package cmd

import (
	"log"
)

func Main() {
	err := Execute()
	if err != nil {
		log.Fatal(err)
	}
}

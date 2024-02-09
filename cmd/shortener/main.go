package main

import (
	"log"

	"github.com/chemax/url-shorter/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Panic(err)
	}

}

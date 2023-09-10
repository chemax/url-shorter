package main

import (
	"github.com/chemax/url-shorter/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Panic(err)
	}

}

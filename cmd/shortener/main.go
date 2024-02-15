package main

import (
	"fmt"
	"log"

	"github.com/chemax/url-shorter/internal/app"
)

var (
	buildVersion = "N\\A"
	buildDate    = "N\\A"
	buildCommit  = "N\\A"
)

func main() {
	printBuildData()
	if err := app.Run(); err != nil {
		log.Panic(err)
	}
}

func printBuildData() {
	fmt.Printf("%s\n%s\n%s\n", buildVersion, buildDate, buildCommit)
}

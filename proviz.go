package main

import (
	"github.com/ebuildy/elasticsearch-proviz/pkg/app"
	log "github.com/sirupsen/logrus"
	"os"
)


func main() {

	log.SetOutput(os.Stdout)

	app := app.New()

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
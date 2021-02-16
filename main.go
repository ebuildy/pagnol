package main

import (
	"github.com/ebuildy/pagnol/app"
	"github.com/ebuildy/pagnol/engines"
	log "github.com/sirupsen/logrus"
	"os"
)

type engine interface {
	Support(action app.ActionItem) bool
	Run(baseConnection app.DefinitionConnection, action app.ActionItem) bool
}

func main() {

	log.SetOutput(os.Stdout)

	app := app.New(func(p app.Parameters, definition app.Definition) {
		enginesImpl := []engine{engines.Elasticsearch(p), engines.HTTP(p)}

		getEngineForAction := func (action app.ActionItem) engine {
			for _, engine := range enginesImpl {
				if engine.Support(action) {
					return engine
				}
			}
			return nil
		}

		for _, action := range definition.Actions {
			engine := getEngineForAction(action)

			engine.Run(definition.Connection, action)
		}
	})

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
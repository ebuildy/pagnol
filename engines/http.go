package engines

import (
	"fmt"
	"github.com/ebuildy/pagnol/app"
	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"strings"
)

type HTTPEngine struct {
	client  *resty.Client
	cli     app.Parameters
}

type ActionSpec struct {
	Method, URL string
}

func HTTP(p app.Parameters) *HTTPEngine {
	client := resty.New()

	if p.Verbose {
		client.SetDebug(true)
	}

	return &HTTPEngine{
		client:  client,
		cli:     p,
	}
}

func (engine *HTTPEngine) Support(action app.ActionItem) bool {
	return action.Kind == "http" || action.Kind == "https"
}

func (engine *HTTPEngine) Run(baseConnection app.DefinitionConnection, action app.ActionItem) bool {
	client := engine.client
	app := engine.cli
	URL := baseConnection.URL

	httpSpec := ActionSpec{}

	mapstructure.Decode(action.Spec, &httpSpec)

	fullURL := fmt.Sprintf("%s%s", URL, httpSpec.URL)

	resp, err := client.R().Execute(strings.ToUpper(httpSpec.Method), fullURL)

	if err != nil {
		app.HandleError(err)
	}

	if resp.IsError() {
		log.Error("error [%s] %s", resp.Status(), resp.Body())
	}

	return resp.IsSuccess()
}

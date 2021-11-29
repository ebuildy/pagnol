package engines

import (
	"crypto/tls"
	"fmt"
	"github.com/ebuildy/pagnol/app"
	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
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

	if len(p.Target.AuthUsername) > 0 {
		client.SetBasicAuth(p.Target.AuthUsername, p.Target.AuthPassword)
	}

	if p.Target.TLSNoVerify {
		client.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })
	} else {
		if len(p.Target.TLSCertificate) > 0 {
			client.SetRootCertificate(p.Target.TLSCertificate)
		}
	}

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

func (engine *HTTPEngine) Run(action app.ActionItem) bool {
	client := engine.client
	app := engine.cli
	URL := engine.cli.Target.URL

	httpSpec := ActionSpec{}

	mapstructure.Decode(action.Spec, &httpSpec)

	fullURL := fmt.Sprintf("%s%s", URL, httpSpec.URL)

	resp, err := client.R().Execute(strings.ToUpper(httpSpec.Method), fullURL)

	if err != nil {
		app.HandleError(err)
	}

	if resp.IsError() {
		app.HandleError(fmt.Errorf("HTTP error [%s] %s", resp.Status(), resp.Body()))
	}

	app.HandleEnd(fmt.Sprintf("%s %s", httpSpec.Method, httpSpec.URL), action.Name, resp.IsSuccess())

	return resp.IsSuccess()
}

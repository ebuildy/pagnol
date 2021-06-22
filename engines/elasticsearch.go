package engines

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ebuildy/pagnol/app"
	"github.com/ebuildy/pagnol/yaml_mapstr"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"strings"
)

type ElasticsearchEngine struct {
	client  *resty.Client
	cli     app.Parameters
}

func Elasticsearch(p app.Parameters) *ElasticsearchEngine {

	client := resty.New().
		SetHeader("Content-Type", "application/json")

	if len(p.Target.AuthUsername) > 0 {
		client.SetBasicAuth(strings.Trim(p.Target.AuthUsername, " "), strings.Trim(p.Target.AuthPassword, " "))
	}

	if p.Target.TLSVerify == false {
		client.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })
	} else {
		if len(p.Target.TLSCertificate) > 0 {
			client.SetRootCertificate(p.Target.TLSCertificate)
		}
	}

	if p.Verbose {
		client.SetDebug(true)
	}

	return &ElasticsearchEngine{
		client:  client,
		cli:     p,
	}
}

func (engine *ElasticsearchEngine) Support(action app.ActionItem) bool {
	return strings.HasPrefix(action.Kind, "org.elasticsearch")
}

func (engine *ElasticsearchEngine) Run(action app.ActionItem) bool {
	client := engine.client
	app := engine.cli
	URL := strings.Trim(engine.cli.Target.URL, "/")
	kind := strings.Replace(action.Kind, "org.elasticsearch/", "", 1)

	URLComponent := kindToURLComponents(kind)
	fullURL := fmt.Sprintf("%s/%s/%s", URL, URLComponent, action.Name)
	log.WithField(action.Kind, action.Name).Debug("Processing index template")

	resp, err := client.R().Get(fullURL)

	if resp.IsSuccess() {
		if app.Verbose {
			log.Debug("resource found, deleting")
		}

		resp, err = client.R().
			Delete(fullURL)

		if err != nil {
			app.HandleError(err)
		}

		if resp.IsError() {
			log.Error("error [%s] %s", resp.Status(), resp.Body())
		}
	}

	spec := yaml_mapstr.CleanupInterfaceMap(action.Spec)

	jsonData, err := json.Marshal(spec)

	if err != nil {
		app.HandleError(err)

		return false
	}

	if !json.Valid([]byte(jsonData)) {
		app.HandleError(errors.New("request is not valid JSON"))

		return false
	}

	resp, err = client.R().
		SetBody(jsonData).
		Put(fullURL)

	if err != nil {
		app.HandleError(err)

		app.HandleEnd(action.Kind, action.Name, false)

		return false
	}

	app.HandleSuccess(resp)
	app.HandleEnd(action.Kind, action.Name, resp.IsSuccess())

	return resp.IsSuccess()
}

func kindToURLComponents(k string) string {
	b := map[string]string {
		"index_template" : "_index_template",
		"snapshot_repository" : "_snapshot",
		"slm_policy" : "_slm/policy",
		"slm-policy" : "_slm/policy",
		"livecycle_policy" : "_ilm/policy",
		"livecycle-policy" : "_ilm/policy",
		"ilm" : "_ilm/policy",
	}

	if v, ok := b[k]; ok {
		return v
	} else {
		return k
	}
}
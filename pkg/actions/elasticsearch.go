package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ebuildy/elasticsearch-proviz/pkg/app"
	"github.com/ebuildy/elasticsearch-proviz/pkg/yaml_mapstr"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ActionItem struct {
	Name string
	Spec map[interface{}]interface{}
}

type Actions struct {

	Connection struct {
		URL string
	}

	IndexTemplates []ActionItem `yaml:"index_templates"`
	SnapshotRepositories []ActionItem `yaml:"snapshot_repositories"`
	SLMPolicies []ActionItem `yaml:"slm_policies"`
}

func New(data []byte) {
	actionsData := Actions{}

	err := yaml.Unmarshal(data, &actionsData)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	client := resty.New().
		SetHeader("Content-Type", "application/json")

	if p.Verbose {
		client.SetDebug(true)
	}
}

func ElasticsearchActions(app *app.Parameters, items []ActionItem) {
	for _, action := range items {
		log.WithField("template", action.Name).Debug("Processing index template")

		if app.Verbose {
			log.Debug("Deleting")
		}

		client.R().
			Delete(fmt.Sprintf("%s/_index_template/%s", actions.Connection.URL, action.Name))

		spec := yaml_mapstr.CleanupInterfaceMap(action.Spec)

		jsonData, err := json.Marshal(spec)

		if err != nil {
			app.HandleError(err)

			continue
		}

		if !json.Valid([]byte(jsonData)) {
			app.HandleError(errors.New("request is not valid JSON"))

			continue
		}

		resp, err := client.R().
			SetBody(jsonData).
			Put(fmt.Sprintf("%s/_index_template/%s", actions.Connection.URL, action.Name))

		if err != nil {
			app.HandleError(err)

			app.HandleEnd("index template", action.Name, false)

			continue
		}

		app.HandleSuccess(resp)
		app.HandleEnd("index template", action.Name, resp.IsSuccess())
	}
}

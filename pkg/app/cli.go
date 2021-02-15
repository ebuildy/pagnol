package app

import (
	"fmt"
	"github.com/ebuildy/elasticsearch-proviz/pkg/actions"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Parameters struct {
	Debug, Verbose, FailFast, DryRun bool
	Actions string
}

func New() *cli.App {

	p := Parameters{}

	return &cli.App{
		Name: "proviz",
		Usage: "elasticsearch provisioner",
		Flags: []cli.Flag {
			&cli.BoolFlag{
				Name: "debug",
				Value: false,
				Usage: "if true, log is verbose",
				Destination: &p.Debug,
			},
			&cli.BoolFlag{
				Name: "verbose",
				Value: false,
				Usage: "if true, log is very verbose",
				Destination: &p.Verbose,
			},
			&cli.BoolFlag{
				Name: "dry-run",
				Value: false,
				Usage: "if true, nothing is sent",
				Destination: &p.DryRun,
			},
			&cli.BoolFlag{
				Name: "fail-fast",
				Value: false,
				Usage: "if true, stop at first error",
				Destination: &p.FailFast,
			},
			&cli.StringFlag{
				Name: "actions",
				Aliases: []string{"a"},
				Usage: "YAML actions",
				Required: true,
				Destination: &p.Actions,
			},
		},
		Action: func (c *cli.Context) error {
			if p.Debug || p.Verbose {
				log.SetLevel(log.DebugLevel)
			}

			log.WithField("file", p.Actions).Debug("Loading file")

			data, err :=  ioutil.ReadFile(p.Actions)

			if err != nil {
				log.Fatal(err)
			}



			actions.ElasticsearchActions(actionsData.IndexTemplates)

				return nil
			},
	}
}


func (cli *Parameters) HandleError(err error) {
	if cli.FailFast {
		log.Fatalf("error: %v", err)
	}

	log.Error("error: %v", err)
}

func (cli *Parameters) HandleSuccess(resp *resty.Response) {
	if resp.IsSuccess() {
		log.Debug(fmt.Sprintf("[%s] %s", resp.Status(), resp.Body()))
	} else {
		log.Warn(fmt.Sprintf("[%s] %s", resp.Status(), resp.Body()))
	}
}

func (cli *Parameters) HandleEnd(kind string, name string, success bool) {
	if success {
		_, _ = color.New(color.Bold, color.FgGreen).Printf("✓ %s %s created!\n", kind, name)
	} else {
		_, _ = color.New(color.Bold, color.FgRed).Printf("✗ %s %s not created!\n", kind, name)
	}
}
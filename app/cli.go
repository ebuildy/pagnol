package app

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type TargetConnection struct {
	AuthUsername, AuthPassword, TLSCertificate, URL string
	TLSNoVerify                                     bool
}

type Parameters struct {
	Debug, Verbose, IgnoreError, DryRun bool
	Actions                             string
	Target                              TargetConnection
}

type ActionItem struct {
	Name, Kind string
	Spec map[interface{}]interface{}
}

type Definition []ActionItem

func New(action func(p Parameters, d Definition)) *cli.App {

	p := Parameters{}

	return &cli.App{
		Name: "pagnol",
		Usage: "run HTTP queries and more",
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
				Usage: "if set, log is very verbose",
				Destination: &p.Verbose,
			},
			&cli.BoolFlag{
				Name: "dry-run",
				Value: false,
				Usage: "if set, nothing is sent",
				Destination: &p.DryRun,
			},
			&cli.BoolFlag{
				Name: "ignore-error",
				Value: false,
				Usage: "if set, no stop actions if error occured",
				Destination: &p.IgnoreError,
			},
			&cli.StringFlag{
				Name: "actions",
				Aliases: []string{"a"},
				Usage: "YAML actions",
				Required: true,
				Destination: &p.Actions,
			},
			&cli.StringFlag{
				Name: "username",
				Required: false,
				Destination: &p.Target.AuthUsername,
				EnvVars: []string{"PAGNOL_TARGET_USERNAME"},
			},
			&cli.StringFlag{
				Name: "password",
				Required: false,
				Destination: &p.Target.AuthPassword,
				EnvVars: []string{"PAGNOL_TARGET_PASSWORD"},
			},
			&cli.StringFlag{
				Name: "tls-ca",
				Required: false,
				Destination: &p.Target.TLSCertificate,
				EnvVars: []string{"PAGNOL_TARGET_TLS_CA"},
			},
			&cli.StringFlag{
				Name: "url",
				Required: true,
				Destination: &p.Target.URL,
				EnvVars: []string{"PAGNOL_TARGET_URL"},
			},
			&cli.BoolFlag{
				Name: "tls-no-verify",
				Required: false,
				Destination: &p.Target.TLSNoVerify,
				EnvVars: []string{"PAGNOL_TARGET_TLS_NO_VERIFY"},
				DefaultText: "false",
			},
		},
		Action: func(context *cli.Context) error {
			if p.Debug || p.Verbose {
				log.SetLevel(log.DebugLevel)
			}

			if len(p.Target.TLSCertificate) > 0 {
				if _, err := os.Stat(p.Target.TLSCertificate); errors.Is(err, os.ErrNotExist) {
					log.Fatal(err)
				}
			}

			if len(p.Target.AuthUsername) > 0 {
				p.Target.AuthUsername = strings.Trim(p.Target.AuthUsername, " ")
			}

			if len(p.Target.AuthPassword) > 0 {
				p.Target.AuthPassword = strings.Trim(p.Target.AuthPassword, " ")
			}

			p.Target.URL = strings.Trim(p.Target.URL, "/")
			
			log.WithField("file", p.Actions).Debug("Loading file")

			data, err :=  ioutil.ReadFile(p.Actions)

			if err != nil {
				log.Fatal(err)
			}

			actionsData := Definition{}

			err = yaml.Unmarshal(data, &actionsData)

			if err != nil {
				log.Fatalf("error: %v", err)
			}

			action(p, actionsData)

			return nil
		},
	}
}


func (cli *Parameters) HandleError(err error) {
	if cli.IgnoreError {
		log.Error("error: %v", err)
	} else {
		log.Fatalf("error: %v", err)
	}
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
		_, _ = color.New(color.Bold, color.FgGreen).Printf("✓ %s %s ok !\n", kind, name)
	} else {
		_, _ = color.New(color.Bold, color.FgRed).Printf("✗ %s %s kaboom !\n", kind, name)
	}
}
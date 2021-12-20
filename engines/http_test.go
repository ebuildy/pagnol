package engines

import (
	"github.com/ebuildy/pagnol/app"
	"testing"
)

func TestHTTPGet(t *testing.T) {
	client := HTTP(app.Parameters{
		Debug:       false,
		Verbose:     false,
		IgnoreError: false,
		DryRun:      false,
		Actions:     "",
		Target: app.TargetConnection{
			URL: "http://httpbin.org",
		},
	})

	success := client.Run(app.ActionItem{
		Name: "",
		Kind: "http",
		Spec: map[interface{}]interface{}{
			"method": "get",
			"url":    "/get",
		},
	})

	if !success {
		t.Errorf("HTTP get not working!")
	}
}

func TestHTTPPost(t *testing.T) {
	client := HTTP(app.Parameters{
		Debug:       false,
		Verbose:     false,
		IgnoreError: false,
		DryRun:      false,
		Actions:     "",
		Target: app.TargetConnection{
			URL: "http://httpbin.org",
		},
	})

	success := client.Run(app.ActionItem{
		Name: "",
		Kind: "http",
		Spec: map[interface{}]interface{}{
			"method":      "post",
			"url":         "/post",
			"content":     "hello",
			"contentType": "application/json",
		},
	})

	if !success {
		t.Errorf("HTTP get not working!")
	}
}

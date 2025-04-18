package kcmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/canonical/rt-conf/src/model"
)

const (
	snapdSocket = "/run/snapd.socket"
	confURL     = "http://localhost/v2/snaps/system/conf"
	jsonbody    = `
{
    "system":{
        "kernel":{
            "dangerous-cmdline-append":"%s"
        }
    }
}`
)

type Result struct {
	Msg string `json:"message"`
	Knd string `json:"kind"`
	Val string `json:"value"`
}

type SnapdResponse struct {
	StatusCode int    `json:"status-code"`
	Status     string `json:"status"`
	Change     string `json:"change"`
	Result     Result `json:"result"`
	Type       string `json:"type"`

	WarningCount     int       `json:"warning-count"`
	WarningTimestamp time.Time `json:"warning-timestamp"`

	Maintenance Result `json:"maintenance"`
}

// createTransport returns an HTTP transport that connects over a Unix socket
func createTransport() *http.Transport {
	return &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", snapdSocket)
		},
	}
}

// sendRequest sends a request to Snapd API and returns the response body
func sendRequest(method, url string, payload []byte) (*http.Response, error) {
	client := &http.Client{Transport: createTransport()}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	return resp, nil
}

func UpdateUbuntuCore(cfg *model.InternalConfig) ([]string, error) {

	cmdline := model.ConstructKeyValuePairs(&cfg.Data.KernelCmdline)
	if len(cmdline) == 0 {
		return nil, fmt.Errorf("no parameters to inject")
	}
	kcmds := model.ParamsToCmdline(cmdline)

	b := []byte(fmt.Sprintf(jsonbody, kcmds))
	resp, err := sendRequest("PUT", confURL, b)
	if err != nil {
		return nil, fmt.Errorf("error communicating with snapd: %s", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var snapResp SnapdResponse
	if err := json.Unmarshal(body, &snapResp); err != nil {
		return nil, fmt.Errorf("error parsing snapd response: %s", err)
	}

	if snapResp.StatusCode >= 400 {
		return nil, fmt.Errorf("snapd error: %s, %s", snapResp.Status,
			snapResp.Result.Msg)
	}

	log.Println("Appended kernel cmdline: ", kcmds)

	return UbuntuCoreConclusion(), nil
}

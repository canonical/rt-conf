package kcmd

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/canonical/rt-conf/src/model"
)

func TestUpdateUbuntuCore(t *testing.T) {
	tests := []struct {
		name     string
		cfg      model.InternalConfig
		mockResp *http.Response
		mockErr  error
		body     string
		expected string // expected error or message
		err      error
	}{
		{
			name: "Empty kernel cmdline",
			cfg: model.InternalConfig{
				Data: model.Config{
					KernelCmdline: model.KernelCmdline{},
				},
			},
			err: errors.New("no parameters to inject"),
		},
		{
			name: "Snapd request fails",
			cfg: model.InternalConfig{
				Data: model.Config{
					KernelCmdline: model.KernelCmdline{"isolcpus=1-3"},
				},
			},
			mockErr: errors.New("connection refused"),
			err:     errors.New("error communicating with snapd: connection refused"),
		},
		{
			name: "Snapd response is invalid JSON",
			cfg: model.InternalConfig{
				Data: model.Config{
					KernelCmdline: model.KernelCmdline{"isolcpus=1-3"},
				},
			},
			mockResp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("not json")),
			},
			err: errors.New("error parsing snapd response"),
		},
		{
			name: "Snapd returns error status",
			cfg: model.InternalConfig{
				Data: model.Config{
					KernelCmdline: model.KernelCmdline{"isolcpus=1-3"},
				},
			},
			mockResp: &http.Response{
				StatusCode: 400,
				Body: io.NopCloser(strings.NewReader(`
				{
					"status": "Bad Request",
					"status-code": 400,
					"result": { "message": "invalid input" }
				}`)),
			},
			err: errors.New("snapd error: Bad Request, invalid input"),
		},
		{
			name: "Success",
			cfg: model.InternalConfig{
				Data: model.Config{
					KernelCmdline: model.KernelCmdline{"isolcpus=1-3"},
				},
			},
			mockResp: &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(strings.NewReader(`
				{
					"status": "Success",
					"status-code": 200,
					"result": { "message": "done" }
				}`)),
			},
			expected: "Please reboot your system to apply the changes.",
			err:      nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sendRequest = func(_, _ string, _ []byte) (*http.Response, error) {
				if tc.mockErr != nil {
					return nil, tc.mockErr
				}
				return tc.mockResp, nil
			}

			msgs, err := UpdateUbuntuCore(&tc.cfg)

			if tc.err != nil {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tc.err)
				}
				if !strings.Contains(err.Error(), tc.err.Error()) {
					t.Fatalf("expected error %q, got %v", tc.err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			found := false
			for _, msg := range msgs {
				if strings.Contains(msg, tc.expected) {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("expected output to contain %q, got %v", tc.expected, msgs)
			}
		})
	}
}

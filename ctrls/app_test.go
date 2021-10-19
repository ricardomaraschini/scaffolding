package ctrls

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestWithAppCtrlBindAddress(t *testing.T) {
	for _, tt := range []struct {
		name string
		addr string
	}{
		{
			name: "happy path",
			addr: ":8181",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := NewAppCtrl(WithAppCtrlBindAddress(tt.addr))
			if ctrl.bind != tt.addr {
				t.Errorf("bind address not set")
			}
		})
	}
}

func TestAppCtrlServeHTTP(t *testing.T) {
	for _, tt := range []struct {
		name    string
		path    string
		method  string
		body    []byte
		experr  string
		expbody []byte
	}{
		{
			name:    "happy path",
			path:    "/",
			method:  http.MethodGet,
			expbody: []byte("OK"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := NewAppCtrl()
			ts := httptest.NewServer(ctrl)
			defer ts.Close()

			requrl := fmt.Sprintf("%s%s", ts.URL, tt.path)
			req, err := http.NewRequest(tt.method, requrl, bytes.NewBuffer(tt.body))
			if err != nil {
				t.Fatalf("unable to create http request: %s", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				if len(tt.experr) == 0 {
					t.Fatalf("unexpected error: %s", err)
				} else if !strings.Contains(err.Error(), tt.experr) {
					t.Fatalf("expected error %q, %q seen", tt.experr, err)
				}
			} else if len(tt.experr) > 0 {
				t.Fatalf("expecting error %q, nil received instead", tt.experr)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("unexpected error reading response body: %s", err)
			}

			if reflect.DeepEqual(body, tt.expbody) {
				return
			}

			t.Errorf(
				"unexpected response %q, expected %q",
				string(body),
				string(tt.expbody),
			)
		})
	}

}

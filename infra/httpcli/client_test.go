package httpcli

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUnmarshalObject(t *testing.T) {
	httpcli := New()
	type Person struct {
		Name    string
		Surname string
	}

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"name": "a", "surname": "b"}`))
		}),
	)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatalf("unexected error creating request: %s", err)
	}

	var out Person
	if err := httpcli.UnmarshalDo(req, &out); err != nil {
		t.Fatalf("error executing request: %s", err)
	}

	exp := Person{Name: "a", Surname: "b"}
	if !reflect.DeepEqual(exp, out) {
		t.Errorf("expected %+v, received %+v", exp, out)
	}
}

func TestUnmarshalSlice(t *testing.T) {
	httpcli := New()
	type Struct struct {
		Test    string
		Another string
	}

	body := []byte(`[
		{"test": "a", "another": "b"},
		{"test": "c", "another": "d"}
	]`)
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(body)
		}),
	)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatalf("unexected error creating request: %s", err)
	}

	var out []Struct
	if err := httpcli.UnmarshalDo(req, &out); err != nil {
		t.Fatalf("error executing request: %s", err)
	}

	exp := []Struct{
		{Test: "a", Another: "b"},
		{Test: "c", Another: "d"},
	}
	if !reflect.DeepEqual(exp, out) {
		t.Errorf("expected %+v, received %+v", exp, out)
	}
}

func TestUnmarshalFail(t *testing.T) {
	httpcli := New()
	body := []byte(` <!SS##$$$#$%$@%$@%$@#$%@#$%@#$%@#@#%@#$%@#ADFAFA#@$$#!% `)
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(body)
		}),
	)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatalf("unexected error creating request: %s", err)
	}

	out := map[string]interface{}{}
	if err := httpcli.UnmarshalDo(req, &out); err == nil {
		t.Fatalf("expected error parsing body, nil received instead")
	}
}

func TestNilURL(t *testing.T) {
	httpcli := New()
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{}`))
		}),
	)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatalf("unexected error creating request: %s", err)
	}
	req.URL = nil

	out := map[string]string{}
	if err := httpcli.UnmarshalDo(req, &out); err == nil {
		t.Fatalf("expected error with nil URL")
	}
}

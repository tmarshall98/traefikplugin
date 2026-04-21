package traefikplugin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tmarshall98/traefikplugin"
)

func TestDemo(t *testing.T) {
	cfg := traefikplugin.CreateConfig()
	cfg.Headers["X-Host"] = "[[.Host]]"
	cfg.Headers["X-Method"] = "[[.Method]]"
	cfg.Headers["X-URL"] = "[[.URL]]"
	cfg.Headers["X-URL"] = "[[.URL]]"
	cfg.Headers["X-Test-FromHeader"] = `[[ .Header.Get "Cf-Connecting-Ip"]]`
	cfg.Headers["X-Plugin"] = "test"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := traefikplugin.New(ctx, next, cfg, "demo-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cf-Connecting-Ip", "TEST")

	handler.ServeHTTP(recorder, req)

	assertHeader(t, req, "X-Host", "localhost")
	assertHeader(t, req, "X-URL", "http://localhost")
	assertHeader(t, req, "X-Method", "GET")
	assertHeader(t, req, "X-Plugin", "test")
	assertHeader(t, req, "X-Test-FromHeader", "TEST")
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()
	got := req.Header.Get(key)
	if got != expected {
		t.Errorf("invalid header value: got=%s, expected=%s", got, req.Header.Get(key))
	}
}

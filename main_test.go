package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	if (w.Code != http.StatusOK) {
		t.Fatalf(`Expected HTTP 200 OK from GET /ping, received %d instead`, w.Code);
	}

	if (w.Body.String() != "pong") {
		t.Fatalf(`Expected response body to equal /"pong/", received %s instead`, w.Body.String());
	}
}
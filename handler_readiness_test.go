package main

import (
	"io"
	"net/http/httptest"
	"testing"
)

func TestHandlerReadiness(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/v1/healthz", nil)

	handlerReadiness(w, r)

	want := struct {
		code     int
		response string
	}{code: 200, response: `{"Status":"OK"}`}

	t.Run("Readiness Handler Test", func(t *testing.T) {
		resp, err := io.ReadAll(w.Body)

		if err != nil {
			t.Fatalf("Could not read response body: %q", err)
		}

		if w.Code != want.code {
			t.Errorf("Incorrect status code, got %d want %d", w.Code, want.code)
		}

		if got := string(resp); got != want.response {
			t.Errorf("Incorrect response body, got %s want %s", got, want.response)
		}
	})
}

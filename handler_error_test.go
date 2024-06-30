package main

import (
	"io"
	"net/http/httptest"
	"testing"
)

func TestHandlerError(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/v1/err", nil)

	handlerError(w, r)

	want := struct {
		code     int
		response string
	}{code: 500, response: `{"error":"Internal Server Error"}`}

	t.Run("Error Handler Test", func(t *testing.T) {
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

package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBodySizeLimit(t *testing.T) {
	payload := bytes.Repeat([]byte("big"), 6*1024*1024)

	req, err := http.NewRequest(http.MethodPost, "/server/run", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Run)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := "reading body\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

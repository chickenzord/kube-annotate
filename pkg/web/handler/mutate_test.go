package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMutateHandlerEmptyBody(t *testing.T) {
	req, err := http.NewRequest("POST", "/mutate", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MutateHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestMutateHandlerWrongContentType(t *testing.T) {
	body := bytes.NewBufferString("{}")
	req, err := http.NewRequest("POST", "/mutate", body)
	req.Header.Set("Content-Type", "text/plain")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MutateHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

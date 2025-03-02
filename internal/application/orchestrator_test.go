package application_test

import (
	"bytes"
	"encoding/json"
	"github.com/KadimovRus/calc_go/internal/application"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateHandler(t *testing.T) {
	orchestrator := application.NewOrchestrator()
	reqBody := []byte(`{"expression":"1+2"}`)
	req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rw := httptest.NewRecorder()
	handler := http.HandlerFunc(orchestrator.CalculateHandler)
	handler.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusCreated, rw.Code)

	var resp map[string]string
	err = json.Unmarshal(rw.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp, "id")
}

func TestExpressionsHandler(t *testing.T) {
	orchestrator := application.NewOrchestrator()
	req, err := http.NewRequest("GET", "/api/v1/expressions", nil)
	assert.NoError(t, err)

	rw := httptest.NewRecorder()
	handler := http.HandlerFunc(orchestrator.ExpressionsHandler)
	handler.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusOK, rw.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rw.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp, "expressions")
}

func TestGetTaskHandler_NoTasks(t *testing.T) {
	orchestrator := application.NewOrchestrator()
	req, err := http.NewRequest("GET", "/internal/task", nil)
	assert.NoError(t, err)

	rw := httptest.NewRecorder()
	handler := http.HandlerFunc(orchestrator.GetTaskHandler)
	handler.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusNotFound, rw.Code)
}

func TestPostTaskHandler_InvalidBody(t *testing.T) {
	orchestrator := application.NewOrchestrator()
	reqBody := []byte(`{"id":"","result":42}`)
	req, err := http.NewRequest("POST", "/internal/task", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rw := httptest.NewRecorder()
	handler := http.HandlerFunc(orchestrator.PostTaskHandler)
	handler.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rw.Code)
}

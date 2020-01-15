package rest_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eroatta/freqtable/adapter/rest"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/frequency-tables", rest.PostFrequencyTable)

	return router
}

func TestPOST_OnFrequencyTableCreationHandler_WithoutBody_ShouldReturnHTTP400(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/frequency-tables", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		assert.FailNow(t, fmt.Sprintf("unexpected unmarshalling err: %v", err))
	}
	assert.Equal(t, "validation_error", response["name"])
	assert.Equal(t, "missing or invalid data", response["message"])
	assert.Equal(t, "invalid request", response["details"].([]interface{})[0].(string))
}

func TestPOST_OnFrequencyTableCreationHandler_WithEmptyBody_ShouldReturnHTTP400(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	body := `{}`
	req, _ := http.NewRequest("POST", "/frequency-tables", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		assert.FailNow(t, fmt.Sprintf("unexpected unmarshalling err: %v", err))
	}
	assert.Equal(t, "validation_error", response["name"])
	assert.Equal(t, "missing or invalid data", response["message"])
	assert.Equal(t, "Key: 'postFrequencyTableCommand.Repository' Error:Field validation for 'Repository' failed on the 'required' tag", response["details"].([]interface{})[0].(string))
}

func TestPOST_OnFrequencyTableCreationHandler_WithWrongDataType_ShouldReturnHTTP400(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	body := `{
		"repository": 1
	}`
	req, _ := http.NewRequest("POST", "/frequency-tables", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		assert.FailNow(t, fmt.Sprintf("unexpected unmarshalling err: %v", err))
	}
	assert.Equal(t, "validation_error", response["name"])
	assert.Equal(t, "missing or invalid data", response["message"])
	assert.Equal(t, "json: cannot unmarshal number into Go struct field postFrequencyTableCommand.repository of type string", response["details"].([]interface{})[0].(string))
}

func TestPOST_OnFrequencyTableCreationHandler_WithInvalidRepository_ShouldReturnHTTP400(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	body := `{
		"repository": "https://github...com/eroatta/freqtable"
	}`
	req, _ := http.NewRequest("POST", "/frequency-tables", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		assert.FailNow(t, fmt.Sprintf("unexpected unmarshalling err: %v", err))
	}
	assert.Equal(t, "validation_error", response["name"])
	assert.Equal(t, "missing or invalid data", response["message"])
	assert.Equal(t, "json: cannot unmarshal number into Go struct field postFrequencyTableCommand.repository of type string", response["details"].([]interface{})[0].(string))
}

// POST with invalid github URL should return 400 Bad Request
// POST with existing github URL should return 400 Bad Request
// POST with valid parameters but a failure while processing should return 500 Internal Error
// POST with valid parameters and successful execution should return 201 Created and the FT info

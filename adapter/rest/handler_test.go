package rest_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

// POST with empty body or any other non required field should return 400 Bad Request
// POST with wrong data type should return 400 Bad Request
// POST with invalid github URL should return 400 Bad Request
// POST with existing github URL should return 400 Bad Request
// POST with valid parameters but a failure while processing should return 500 Internal Error
// POST with valid parameters and successful execution should return 201 Created and the FT info

package rest_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/eroatta/freqtable/adapter/rest"
	"github.com/eroatta/freqtable/entity"
	"github.com/stretchr/testify/assert"
)

func TestPOST_OnFrequencyTableCreationHandler_WithoutBody_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer(nil)

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
	router := rest.NewServer(nil)

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
	assert.Equal(t, "invalid field 'repository' with value null or empty", response["details"].([]interface{})[0].(string))
}

func TestPOST_OnFrequencyTableCreationHandler_WithWrongDataType_ShouldReturnHTTP400(t *testing.T) {
	router := rest.NewServer(nil)

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
	router := rest.NewServer(nil)

	w := httptest.NewRecorder()
	body := `{
		"repository": "./github.com/eroatta/freqtable"
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
	assert.Equal(t, "invalid field 'repository' with value ./github.com/eroatta/freqtable", response["details"].([]interface{})[0].(string))
}

func TestPOST_OnFrequencyTableCreationHandler_WithInternalError_ShouldReturnHTTP500(t *testing.T) {
	router := rest.NewServer(mockUsecase{
		ft:  entity.FrequencyTable{},
		err: errors.New("error cloning repository http://github.com/eroatta/freqtable"),
	})

	w := httptest.NewRecorder()
	body := `{
		"repository": "http://github.com/eroatta/freqtable"
	}`
	req, _ := http.NewRequest("POST", "/frequency-tables", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		assert.FailNow(t, fmt.Sprintf("unexpected unmarshalling err: %v", err))
	}
	assert.Equal(t, "internal_error", response["name"])
	assert.Equal(t, "internal server error", response["message"])
	assert.Equal(t, "error cloning repository http://github.com/eroatta/freqtable", response["details"].([]interface{})[0].(string))
}

func TestPOST_OnFrequencyTableCreationHandler_WithSuccess_ShouldReturnHTTP201(t *testing.T) {
	now := time.Now()
	ft := entity.FrequencyTable{
		ID:          int64(123112312),
		Name:        "http://github.com/eroatta/freqtable",
		DateCreated: now,
		LastUpdated: now,
	}

	router := rest.NewServer(mockUsecase{
		ft:  ft,
		err: nil,
	})

	w := httptest.NewRecorder()
	body := `{
		"repository": "http://github.com/eroatta/freqtable"
	}`
	req, _ := http.NewRequest("POST", "/frequency-tables", strings.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		assert.FailNow(t, fmt.Sprintf("unexpected unmarshalling err: %v", err))
	}
	responseId, _ := strconv.Atoi(fmt.Sprintf("%.0f", response["id"]))
	assert.Equal(t, 123112312, responseId)
	assert.Equal(t, "http://github.com/eroatta/freqtable", response["name"])
	assert.Equal(t, now.Format(time.RFC3339), response["date_created"])
	assert.Equal(t, now.Format(time.RFC3339), response["last_updated"])
}

type mockUsecase struct {
	ft  entity.FrequencyTable
	err error
}

func (m mockUsecase) Create(ctx context.Context, url string) (entity.FrequencyTable, error) {
	return m.ft, m.err
}

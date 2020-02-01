package rest

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/eroatta/freqtable/usecase"
	"github.com/eroatta/token/conserv"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"
)

// requestValidator represents a validator capable of analyzing the values of the incoming
// request bodies.
var requestValidator = validator.New()

// NewServer creates a new gingonic Engine that handles HTTP requests.
func NewServer(ftUsecase usecase.CreateFrequencyTableUsecase) *gin.Engine {
	internal := server{
		createFreqTableUseCase: ftUsecase,
	}

	r := gin.Default()
	r.GET("/ping", pingHandler)
	r.POST("/frequency-tables", internal.postFrequencyTable)

	return r
}

type server struct {
	createFreqTableUseCase usecase.CreateFrequencyTableUsecase
}

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

type postFrequencyTableCommand struct {
	Repository string `json:"repository" validate:"url"`
}

type freqTableResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DateCreated string `json:"date_created"`
	LastUpdated string `json:"last_updated,omitempty"`
}

type errorResponse struct {
	Name    string   `json:"name"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

func (s server) postFrequencyTable(ctx *gin.Context) {
	var cmd postFrequencyTableCommand

	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		log.WithError(err).Debug("failed to bind JSON body")
		setBadRequestOnBindingResponse(ctx, err)
		return
	}

	if err := requestValidator.Struct(cmd); err != nil {
		log.WithError(err).Debug("failed while validating the command")
		setBadRequestOnValidationResponse(ctx, err)
		return
	}

	ft, err := s.createFreqTableUseCase.Create(ctx, cmd.Repository)
	if err != nil {
		log.WithError(err).Error("unexpected error")
		setInternalErrorResponse(ctx, err)
		return
	}

	response := freqTableResponse{
		ID:          ft.ID,
		Name:        ft.Name,
		DateCreated: ft.DateCreated.Format(time.RFC3339),
		LastUpdated: ft.LastUpdated.Format(time.RFC3339),
	}
	ctx.JSON(http.StatusCreated, response)
}

func newBadRequestResponse() errorResponse {
	return errorResponse{
		Name:    "validation_error",
		Message: "missing or invalid data",
		Details: make([]string, 0),
	}
}

func setBadRequestOnBindingResponse(ctx *gin.Context, err error) {
	errResponse := newBadRequestResponse()
	errResponse.Details = append(errResponse.Details, err.Error())

	ctx.JSON(http.StatusBadRequest, errResponse)
}

func setBadRequestOnValidationResponse(ctx *gin.Context, err error) {
	errResponse := newBadRequestResponse()
	for _, err := range err.(validator.ValidationErrors) {
		field := strings.ToLower(strings.Join(conserv.Split(err.Field()), "_"))
		var value interface{}
		if val, ok := err.Value().(string); ok && val == "" {
			value = "null or empty"
		} else {
			value = err.Value()
		}
		errResponse.Details = append(errResponse.Details,
			fmt.Sprintf("invalid field '%s' with value %v", field, value))
	}

	ctx.JSON(http.StatusBadRequest, errResponse)
}

func setInternalErrorResponse(ctx *gin.Context, err error) {
	errResponse := errorResponse{
		Name:    "internal_error",
		Message: "internal server error",
		Details: []string{err.Error()},
	}

	ctx.JSON(http.StatusInternalServerError, errResponse)
}

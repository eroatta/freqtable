package rest

import (
	"fmt"
	"strings"
	"time"

	"github.com/eroatta/freqtable/usecase"
	"github.com/eroatta/token/conserv"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"
)

type Server interface {
	//Run()
	PostFrequencyTable(ctx *gin.Context)
}

func NewServer(ftUsecase usecase.CreateFrequencyTableUsecase) Server {
	return server{
		createFreqTableUseCase: ftUsecase,
	}
}

type server struct {
	createFreqTableUseCase usecase.CreateFrequencyTableUsecase
}

func PingHandler(c *gin.Context) {
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

var validate = validator.New()

func (s server) PostFrequencyTable(ctx *gin.Context) {
	var cmd postFrequencyTableCommand

	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		log.WithError(err).Debug("failed to bind JSON body")
		errResponse := errorResponse{
			Name:    "validation_error",
			Message: "missing or invalid data",
			Details: []string{
				err.Error(),
			},
		}
		ctx.JSON(400, errResponse)
		return
	}

	err := validate.Struct(cmd)
	if err != nil {
		errResponse := errorResponse{
			Name:    "validation_error",
			Message: "missing or invalid data",
			Details: make([]string, 0),
		}
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
		ctx.JSON(400, errResponse)
		return
	}

	ft, err := s.createFreqTableUseCase.Create(ctx, cmd.Repository)
	if err != nil {
		errResponse := errorResponse{
			Name:    "internal_error",
			Message: "internal server error",
			Details: []string{err.Error()},
		}
		ctx.JSON(500, errResponse)
		return
	}

	resp := freqTableResponse{
		ID:          ft.ID,
		Name:        ft.Name,
		DateCreated: ft.DateCreated.Format(time.RFC3339),
		LastUpdated: ft.LastUpdated.Format(time.RFC3339),
	}
	ctx.JSON(201, resp)
}

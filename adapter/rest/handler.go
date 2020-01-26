package rest

import (
	"fmt"
	"strings"
	"time"

	"github.com/eroatta/token/conserv"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"
)

type Server interface {
	Run()
}

func NewServer() Server {
	return nil
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
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	DateCreated time.Time `json:"date_created"`
	LastUpdated time.Time `json:"last_update,omitempty"`
}

type errorResponse struct {
	Name    string   `json:"name"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

var validate = validator.New()

func PostFrequencyTable(ctx *gin.Context) {
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

	resp := freqTableResponse{
		ID:          1,
		Name:        cmd.Repository,
		DateCreated: time.Now(),
		LastUpdated: time.Time{},
	}
	ctx.JSON(201, resp)
}

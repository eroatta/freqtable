package rest

import (
	"time"

	"github.com/gin-gonic/gin"
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
	Repository string `json:"repository" binding:"required"`
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

func PostFrequencyTable(ctx *gin.Context) {
	var cmd postFrequencyTableCommand

	if err := ctx.BindJSON(&cmd); err != nil {
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

	resp := freqTableResponse{
		ID:          1,
		Name:        cmd.Repository,
		DateCreated: time.Now(),
		LastUpdated: time.Time{},
	}
	ctx.JSON(201, resp)
}

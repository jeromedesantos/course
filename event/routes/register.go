package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jeromedesantos/rest-api/models"
)

func registerForEvent(context *gin.Context) {
	var (
		event *models.Event
		eventId int64
		userId int64
		err error
	)

	userId = context.GetInt64("userId") // --> ambil token userId
	eventId, err = strconv.ParseInt(context.Param("id"), 10, 64) // --> ambil parameter id, convert ke int64
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse event id.",
		})
		return
	}

	// cek apakah event ada
	event, err = models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch event",
		})
		return
	}

	err = event.Register(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not register for event",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Registered for event successfully",
	})
}

  
func cancelRegistration(context *gin.Context) {
	var (
		event models.Event
		eventId int64
		userId int64
		err error
	)

	userId = context.GetInt64("userId") // --> ambil token userId
	eventId, err = strconv.ParseInt(context.Param("id"), 10, 64) // --> ambil parameter id, convert ke int64
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse event id.",
		})
		return
	}

	event.ID = eventId

	err = event.CancelRegistration(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not cancel registration",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Cancelled!",
	})
}
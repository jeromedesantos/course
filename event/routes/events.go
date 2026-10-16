package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jeromedesantos/rest-api/models"
)

// http server :
// --> get request, send back response
func getEvents(context *gin.Context) {
	var (
		events	[]models.Event
		err		error
	)

	events, err = models.GetAllEvents() // --> ambil semua data
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch events. Try again later.",
		}) // gin.H --> map[string]interface{}
		return
	}
	context.JSON(http.StatusOK, events)  // --> http status code 200
}

func getEvent(context *gin.Context) {
	var (
		eventId	int64
		err		error
	)

	// tarik path parameter → convert str to int64
	eventId, err = strconv.ParseInt(context.Param("id"), 10, 64) // --> basis 10 (desimal), 64 bit
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse event id.",
		})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not parse event.",
		})
		return
	}

	context.JSON(http.StatusOK, event)
}

func createEvent(context *gin.Context) {
	var (
		event models.Event // --> wadah sementara buat nampung data
		userId int64
		err error
	)


	err = context.ShouldBindJSON(&event) // --> scanner
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data.",
		})
		return // --> berhentiin fungsi, jangan lanjutin ke bawah
	}

	// set user id from token
	userId = context.GetInt64("userId")
	event.UserID = userId 

	err = event.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message":"Could not create event. Try again later.",
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Event created!.",
		"event": event,
	})
}

func updateEvent(context *gin.Context) {
	var (
		eventId	int64
		updatedEvent models.Event
		userId int64
		event *models.Event
		err error
	)

	eventId, err = strconv.ParseInt(context.Param("id"), 10, 64) 
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse event id.",
		})
		return
	}

	userId = context.GetInt64("userId") // tarik user id dari token
	event, err = models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch the event.",
		})
		return
	}
	
	// cek apakah user yang login adalah owner dari event yang mau diupdate
	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "You are not authorized to update event.",
		})
		return
	}

	err = context.ShouldBindJSON(&updatedEvent)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data.",
		})
		return
	}

	updatedEvent.ID = eventId
	err = updatedEvent.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not update event. Try again later.",
		})
		return
	}
	
	context.JSON(http.StatusOK, gin.H{
		"message": "Event updated successfully.",
		"event": updatedEvent,
	})
}

func deleteEvent(context *gin.Context) {
	var (
		eventId	int64
		event *models.Event
		err error
	)

	eventId, err = strconv.ParseInt(context.Param("id"), 10, 64) 
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse event id.",
		})
		return
	}

	event, err = models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch event.",
		})
		return
	}

	err = event.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not delete event.",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Event deleted successfully.",
	})
}


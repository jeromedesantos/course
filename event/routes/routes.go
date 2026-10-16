package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jeromedesantos/rest-api/middlewares"
)

func RegisterRoutes(server *gin.Engine) {
	var (
		authenticated *gin.RouterGroup
	)

	// GET, POST, PUT, PATCH, DELETE
	server.GET("/events", getEvents) // --> path, handler
	server.GET("/events/:id", getEvent) // /events/1, /events/2 etc

	// grouped routes
	authenticated = server.Group("/")
	authenticated.Use(middlewares.Authenticate) // --> applies to all routes inside this group

	authenticated.POST("/events", createEvent)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.DELETE("/events/:id", deleteEvent)

	authenticated.POST("/events/:id/register", registerForEvent)
	authenticated.DELETE("/events/:id/register", cancelRegistration)

	server.POST("/signup", signup)
	server.POST("/login", login)
}
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jeromedesantos/rest-api/db"
	"github.com/jeromedesantos/rest-api/routes"
)

func main() {
	db.InitDB() // --> init database (buat table)
	server := gin.Default()

	routes.RegisterRoutes(server) // --> daftarin route
	
	server.Run(":8080") // --> start server 
	// in localhost:8080
}

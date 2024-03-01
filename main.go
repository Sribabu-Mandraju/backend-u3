package main

import (
	"backend/database"
	routes "backend/routes"
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func main() {
	client = database.DBinstance()
	defer client.Disconnect(context.Background())
	router := gin.Default()
	config := cors.DefaultConfig()
	
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	router.Use(cors.New(config))

	//routes
	routes.AdminRoutes(router)
	routes.AuthRoutes(router)
	routes.AdminAuthRoutes(router)
	

	router.Use(gin.Logger())
	port := "8080"
	router.Run(":" + port)
}
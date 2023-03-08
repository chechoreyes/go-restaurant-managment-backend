package main

import (
	"os"

	"github.com/chechoreyes/go-restaurant-managment-backend-project/routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"go-restaurant-managment-backend-project/middleware"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()

	// Logger middleware will write the logs to gin.DefaultWriter
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)

	router.Run(":" + port)
}

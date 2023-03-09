package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/chechoreyes/go-restaurant-managment-backend-project/database"
	"github.com/chechoreyes/go-restaurant-managment-backend-project/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Get food_id property of request
		foodId := c.Param("food_id")
		var food models.Food

		// Search the food_id in the collection and put in the var food address memory
		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)

		if err != nil {
			// Send response with StatusInternalServerError code and the JSON with de error written
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error ocurred while fetching the food item"})
		}

		// Send successful response
		c.JSON(http.StatusOK, food)

		defer cancel()
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu
		var food models.Food

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		validationErr := validate.Struct(food)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		}

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
		defer cancel()

		if err != nil {
			msg := fmt.Sprintf("menu was not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		food.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		food.Food_id = food.ID.Hex()
		var num = toFixed(*food.Price, 2)
		food.Price = &num

		result, insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr != nil {
			msg := fmt.Sprintf("Food item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)

	}
}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func round(num float64) int {
	return 1
}

func toFixed(num float64, precision int) float64 {
	return 1.0
}

package controllers

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/chechoreyes/go-restaurant-managment-backend-project/database"
	"github.com/chechoreyes/go-restaurant-managment-backend-project/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// RECORDS PER PAGE
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		// PAGINATION
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		// RECORD STARTS
		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		// AGGREGATION TO SEARCH RESUTLS
		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}

		projectStage := bson.D{
			{
				Key: "$project", Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "total_count", Value: 1},
					{Key: "food_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
				}}}

		result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing food items"})
		}
		var allFoods []bson.M
		if err = result.All(ctx, &allFoods); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allFoods[0])
	}
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
			return
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
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var food models.Food
		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		foodId := c.Param("food_id")

		var updateObj primitive.D

		if food.Name != nil {
			updateObj = append(updateObj, bson.E{Key: "name", Value: food.Name})
		}

		if food.Price != nil {
			updateObj = append(updateObj, bson.E{Key: "price", Value: food.Price})
		}

		if food.Food_image != nil {
			updateObj = append(updateObj, bson.E{Key: "food_image", Value: food.Food_image})
		}

		if food.Menu_id != nil {
			err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
			defer cancel()
			if err != nil {
				msg := fmt.Sprintf("message: Menu was not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "menu", Value: food.Price})
		}

		food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: food.UpdatedAt})

		upsert := true
		filter := bson.M{"food_id": foodId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := foodCollection.UpdateOne(
			ctx, filter, bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprint("Food item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}

		defer cancel()
		c.JSON(http.StatusOK, result)

	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

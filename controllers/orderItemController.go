package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func GetOrderItemByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

// M is an unordered representation of a BSON document (map)
func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {
	return 
}

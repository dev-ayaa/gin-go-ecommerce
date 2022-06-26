package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yusuf/gin-gonic-ecommerce/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAddress() gin.HandlerFunc {

	return nil
}

func EditHomeAddress() gin.HandlerFunc {
	return nil
}

func EditWorkAddress() gin.HandlerFunc {
	return nil
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		var address []model.Address

		userID := c.Query("id")
		if userID == "" {
			log.Println("empty user id")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid search value for user id"})
			return
		}

		//creating a new id for the user
		userAddID, err := primitive.ObjectIDFromHex("user_id")
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, "Error from the server")
			return
		}

		ctx , cancelCtx := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancelCtx() 

		filter := bson.D{primitive.E{Key: "_id", Value: userAddID}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: address}}}}

		_,  err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil{
			c.IndentedJSON(http.StatusBadRequest, "cannot change the address")
			return
		}
		defer cancelCtx()
		ctx.Done()
		c.IndentedJSON(http.StatusOK, "successfully delete the address")
 
	}
}

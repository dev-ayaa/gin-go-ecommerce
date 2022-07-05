package controller

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yusuf/gin-gonic-ecommerce/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {

	return func(c *gin.Context) {
		var address model.Address

		queryUserID := c.Query("id")
		if queryUserID == "" {
			log.Println("invalid user id")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, "invalid user id")
			c.Abort()
			return
		}
		userAdd, err := primitive.ObjectIDFromHex(queryUserID)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "Add Address: Internal server error")
		}

		address.AddressID = primitive.NewObjectID()
		if err = c.BindJSON(&address); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, "Invalid address")
		}

		ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancelCtx()

		//the aggregate function

		//match the user id with the new hax value
		matchAddress := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: userAdd}}}}

		//unwind the address
		unWindAddress := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}

		//group
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		addCursor , err :=UserCollection.Aggregate(ctx, mongo.Pipeline{matchAddress, unWindAddress, group})
		if err != nil{
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "cannot add address to the database")
			return
		}

	
		var newAddress []bson.D

		if err = addCursor.All(ctx, &newAddress); err != nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		for _, jAdd := range newAddress {
			c.IndentedJSON(http.StatusOK,jAdd["count"])
			c.IndentedJSON(http.StatusOK, &newAddress)
		}

	}
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
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Error from the server")
			return
		}

		ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelCtx()

		filter := bson.D{primitive.E{Key: "_id", Value: userAddID}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: address}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "cannot change the address")
			return
		}
		defer cancelCtx()
		ctx.Done()
		c.IndentedJSON(http.StatusOK, "successfully delete the address")

	}
}

package controller

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yusuf/gin-gonic-ecommerce/database"
	"github.com/yusuf/gin-gonic-ecommerce/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	userCollection    *mongo.Collection
	productCollection *mongo.Collection
}

func NewApplication(userCollection, productCollection *mongo.Collection) *Application {
	return &Application{
		userCollection:    userCollection,
		productCollection: productCollection,
	}
}

func (app *Application) AddProductToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get the user id
		//get the product id
		//use primitive to check for a unique product id
		//Add item to the user cart

		//query to get the product id
		queryProductID := c.Query("product_id")
		if queryProductID == "" {
			log.Println("No product id")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product has no id"))
			return
		}

		//query to get the user id
		queryUserID := c.Query("user_id")
		if queryUserID == " " {
			log.Println("No user id")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user has no id"))
			return
		}

		//if the product id and the user id is not empty and valid
		productID, err := primitive.ObjectIDFromHex(queryProductID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelCtx()

		err = database.AddProductToCart(ctx, app.userCollection, app.productCollection, productID, queryUserID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(http.StatusOK, "successfully added product to cart")

		// userID, err := primitive.ObjectIDFromHex(queryUserID)

	}
}

func (app *Application) RemoveItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get the user id
		//get the product id as well
		//get the product id from the database
		//remove the item frpm the database of the user of that id

		queryProductID := c.Query("product_id")
		if queryProductID == "" {
			log.Println("product id not in database")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id not in database"))
			return
		}

		queryUserID := c.Query("user_id")
		if queryUserID == "" {
			log.Println("user id not in the database")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id not in the database"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(queryProductID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelCtx()

		err = database.RemoveCartItem(ctx, app.userCollection, app.productCollection, queryProductID, productID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(http.StatusOK, "Successfully remove item from cart")

	}
}

func (app *Application) GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		//to get items from the cart i need to know which user in particular and what
		//item does the user have
		//using aggregate function --MATCH -- UNMATCH -GROUP
		var cartItem model.User

		//query to get the id of the user
		queryUserID := c.Query("id")

		//check for valid id
		if queryUserID == "" {
			log.Println("Empty user id")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid user id: empty user id"})
			c.Abort()
			return

		}

		//if the id is valid change to hex value
		userID, err := primitive.ObjectIDFromHex(queryUserID)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancelCtx()

		err = UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userID}}).Decode(&cartItem)

		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "Item not found in cart")
			return
		}

		//match the user id with all available user data
		filterMatch := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: userID}}}}

		//find the cart of the user and ungroup it i.e to have access to all the cart
		unWind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$user_cart"}}}}

		//group the total value of the item from the cart
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "user_cart.price"}}}}}}

		//aggregate function
		aggResult, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filterMatch, unWind, group})
		if err != nil {
			log.Println(err)
			// c.IndentedJSON(http.StatusInternalServerError, "invalid aggregate for item")
			return
		}

		var itemList []bson.D

		if err = aggResult.All(ctx, &itemList); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		for _, jItem := range itemList {
			c.IndentedJSON(http.StatusOK, jItem["total"])
			c.IndentedJSON(http.StatusOK, cartItem.UserCart)
		}
		ctx.Done()

	}
}

func (app *Application) BuyItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		queryUserID := c.Query("user_id")
		if queryUserID == "" {
			log.Print("no user id , id not available")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("no user id , id not available"))
			return
		}

		ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelCtx()

		err := database.BuyItemFromCart(ctx, app.userCollection, app.productCollection, queryUserID)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(http.StatusOK, "successfully buy item from cart")
	}

}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get the product id
		//get the right user id that want to buy the item
		//get the amount paid for the item
		queryProductID := c.Query("product_id")
		if queryProductID == "" {
			log.Println("no product id")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("no product id"))
			return
		}
		queryUserID := c.Query("user_id")
		if queryProductID == "" {
			log.Println("no user id")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("no user id"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(queryProductID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelCtx()

		err = database.InstantBuy(ctx, app.userCollection, app.productCollection, productID, queryUserID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(http.StatusOK, "successfully buy the item")

	}

}

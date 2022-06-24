package controller

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yusuf/gin-gonic-ecommerce/database"
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
			c.AbortWithError(http.StatusBadRequest, errors.New("product id not in database"))
			return
		}

		queryUserID := c.Query("user_id")
		if queryUserID == "" {
			log.Println("user id not in the database")
			c.AbortWithError(http.StatusBadRequest, errors.New("user id not in the database"))
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

func (app *Application) GetItemFromCart() *gin.HandlerFunc {
	return nil
}

func (app *Application) BuyItemFromCart() *gin.HandlerFunc {
	return nil
}

func (app *Application) InstantBuy() *gin.HandlerFunc {
	return nil
}

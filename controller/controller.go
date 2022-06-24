package controller

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yusuf/gin-gonic-ecommerce/model"
	"github.com/yusuf/gin-gonic-ecommerce/token"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//HashPassword to hashed the user typed in password
func HashPassword(password string) string {
	return ""
}

//VerifyPassword to Verify the user password  with respect to the hashed password
func VerifyPassword(userPassword, hashedPassword string) (bool, string) {

	return true, ""
}

//SignUp allow the user to signup handlers
func SignUp() gin.HandlerFunc {

	return func(c *gin.Context) {
		var user model.User

		ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancelCtx()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//checking validation for the user models
		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
		}

		//checking for used email
		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("User already existed")})
		}

		//checking for used phone number
		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})

		defer cancelCtx()

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("this phone number is already used")})
		}

		password := HashPassword(*user.Password)
		token, refresh_token, _ := generate.TokenGenerator(*user.FirstName, *user.LastName, *user.Email)

		user.Password = &password
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()
		user.CreatedAt, _ = time.Parse("RFC3339", time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse("RFC3339", time.Now().Format(time.RFC3339))
		user.AddressDetails = make([]model.Address, 0)
		user.OrderStatus = make([]model.Order, 0)
		user.UserCart = make([]model.ProductUser, 0)
		user.Token = &token
		user.RefreshToken = &refresh_token
		
		_, insertErr := UserCollection.InsertOne(ctx, user)
		if insertErr != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"user cannot be created or insert into database"})
			return
		}
		defer cancelCtx()

		c.JSON(http.StatusCreated, "Signed in successfully")

	}
}

//Login handlers
func Login() gin.HandlerFunc {

	return func(c *gin.Context){

		ctx, cancelCtx := context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancelCtx()


	}
}

//AdminAddProduct
func AdminAddProduct() *gin.HandlerFunc {
	return nil
}

func SearchProduct() *gin.HandlerFunc {
	return nil
}

func SearchProductByQuery() *gin.HandlerFunc {
	return nil
}

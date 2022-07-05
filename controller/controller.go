package controller

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/yusuf/gin-gonic-ecommerce/database"
	"github.com/yusuf/gin-gonic-ecommerce/model"
	"github.com/yusuf/gin-gonic-ecommerce/token"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

//database collections
var UserCollection *mongo.Collection = database.UserData(database.Client, "User")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Product")

//Validator for models
var Validate = validator.New()

//HashPassword to hashed the user typed in password in the database
func HashPassword(password string) string {
	bytvalue, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		log.Println(err)
		log.Panicln(err)
		return ""
	}
	return string(bytvalue)
}

//VerifyPassword to Verify the user password  with respect to the hashed password
func VerifyPassword(userPassword, hashedPassword string) (bool, string) {
	valid := true
	msg := ""
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		valid = false
		msg = "invalid Login detail:: Incorrect password"
		//  log.Println()
		return valid, msg
	}
	return valid, msg
}

//SignUp allow the user to signup handlers
func SignUp() gin.HandlerFunc {

	return func(c *gin.Context) {
		var user model.User

		ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancelCtx()

		//Marshalling and UnMarshalling the user struct model
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//checking validation for the user models
		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
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
		user.ID = primitive.NewObjectID() //new id since the user will just be signing up
		user.UserID = user.ID.Hex()
		user.CreatedAt, _ = time.Parse("RFC3339", time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse("RFC3339", time.Now().Format(time.RFC3339))
		user.AddressDetails = make([]model.Address, 0)
		user.OrderStatus = make([]model.Order, 0)
		user.UserCart = make([]model.ProductUser, 0)
		user.Token = &token
		user.RefreshToken = &refresh_token

		//Add the user to the UserCollection in the database
		_, insertErr := UserCollection.InsertOne(ctx, user)
		if insertErr != nil {
			log.Println(insertErr)
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user cannot be created or insert into database"})
			return
		}
		defer cancelCtx()

		log.Println("Signed in successfully")
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusCreated, "Signed in successfully")

	}
}

//Login handlers
func Login() gin.HandlerFunc {

	return func(c *gin.Context) {
		var user model.User
		ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancelCtx()

		if err := c.BindJSON(&user); err != nil {
			log.Println(err.Error())
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Checking if the user exist in the database
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancelCtx()
		if err != nil {
			log.Println(err)
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "incorrect login email"})
			return
		}
		validPassword, passwordMsg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancelCtx()
		if !validPassword {
			c.JSON(http.StatusInternalServerError, gin.H{"error": passwordMsg})
			log.Println(passwordMsg)
			return
		}

		token, refresh_token, _ := generate.TokenGenerator(*foundUser.FirstName, *foundUser.LastName, *foundUser)
		defer cancelCtx()

		generate.UpdateAllToken(token, refresh_token)
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusFound, foundUser)
	}
}

//AdminAddProduct
func AdminAddProduct() gin.HandlerFunc {
	return nil
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get all list of the product
		var productList []model.Product
		ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelCtx()

		//pass empty query to get all the product  in the database
		cursor, err := ProductCollection.Find(ctx, bson.D{})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Heh, something went wrong, try again after few minutes")
			return
		}

		//this will iterates through a collection and store all the values
		err = cursor.All(ctx, &productList)
		if err != nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.Err() ; err != nil{
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "invalid result")
			return
		}

		c.IndentedJSON(http.StatusOK, productList)

	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context){
		var productQuery []model.Product

		queryParams := c.Query("name")
		if queryParams == ""{
			log.Println("empty query result")
			c.Header("Content-Type","application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "empty query parameter"})
			c.Abort()
			return
		}

		//if the queryParams is not empty 
		//initialize a context
		ctx, cancelCtx := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancelCtx()

		//pass empty query to get all the product  in the database
		searchProduct, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex" :queryParams}})
		if err != nil{
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "unable to fetch the product from the database")
			return
		}
		err = searchProduct.All(ctx, &productQuery)
		if err != nil{
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "invalid result from databse")
			return
		}

		defer searchProduct.Close(ctx)
		if err = searchProduct.Err(); err != nil{
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "invalid request from the database")
			return
		}
		defer cancelCtx()

		c.IndentedJSON(http.StatusOK, searchProduct)

	}
}

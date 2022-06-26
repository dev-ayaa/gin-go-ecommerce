package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/yusuf/gin-gonic-ecommerce/controller"
	"github.com/yusuf/gin-gonic-ecommerce/database"
	"github.com/yusuf/gin-gonic-ecommerce/middleware"
	"github.com/yusuf/gin-gonic-ecommerce/routes"
)

func main() {

	portNumber := os.Getenv("PORT")

	if portNumber == "" {
		portNumber = "8080"
	}

	app := controller.NewApplication(database.ProductData(database.Client, "Product"), database.UserData(database.Client, "User"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddProductToCart())
	router.GET("/removeitem", app.RemoveItemFromCart())
	router.GET("/cartcheckout", app.BuyItemFromCart())
	router.GET("/getfromcart", app.GetItemFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + portNumber))
}

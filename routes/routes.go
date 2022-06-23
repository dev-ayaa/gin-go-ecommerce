package routes


import(
	"github.com/gin-gonic/gin"
	"github.com/yusuf/ecommerce-cart/controller"
)

//UserRoutes routes for the user 
func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("/user/signup",controller.SignUp())
	incomingRoutes.POST("/user/login", controller.Login()) 
	incomingRoutes.POST("/admin/add-product", controller.AdminAddProduct())
	incomingRoutes.GET("/user/view-product", controller.SearchProduct())
	incomingRoutes.GET("user/search-product",controller.SearchProductByQuery())

}
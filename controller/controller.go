package controller

import "github.com/gin-gonic/gin"


//HashPassword to hashed the user typed in password
func HashPassword(password string) string {
	return ""
}

//VerifyPassword to Verify the user password  with respect to the hashed password
func VerifyPassword(userPassword, hashedPassword string) (bool, string) {

	return true, ""
}

//SignUp handlers
func SignUp() *gin.HandlerFunc {

	return nil
}

//Login handlers
func Login() *gin.HandlerFunc {

	return nil
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


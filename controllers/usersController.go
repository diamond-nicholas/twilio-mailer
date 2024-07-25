package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicholas/go-jwt/initializers"
	"github.com/nicholas/go-jwt/models"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	//get email/pass from req.body
		var body struct {
			Name string
			Email string
			Password string
		}

	if	c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Failed to read body",
		})
		return 
	}
	//hash the password

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Failed to hash password",
		})
		return
	}
	//create the user
	user := models.User{Name: body.Name ,Email: body.Email, Password: string(hash)  }

result := initializers.DB.Create(&user)

if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Failed to create user",
		})
		return
	
}

	//response
	c.JSON(http.StatusOK, gin.H{})
}
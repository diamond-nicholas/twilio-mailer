package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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


err = SendEmail(user.Email, user.Name, "The Emailer", "templates/welcome.html", struct{ Name string }{Name: user.Name})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send welcome email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}


func Login(c *gin.Context){
		var body struct {
			Email string
			Password string
		}

	if	c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Failed to read body",
		})
		return 
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Invalid Email",
		})
		return 
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Invalid  Password",
		})
		return 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub" : user.ID,
		"exp" : time.Now().Add(time.Hour * 24 * 30).Unix(),
	})



	tokenString, err := token.SignedString([]byte(os.Getenv("JWTSECRET")))

	if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Unable to create token",
		})
		return 
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func Validate(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"message" : 	"I'm logged in",
	})
}
package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nicholas/go-jwt/initializers"
	"github.com/nicholas/go-jwt/models"
	"github.com/xuri/excelize/v2"
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

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Failed to hash password",
		})
		return
	}
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

func UploadAndSendEmails(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	xlFile, err := excelize.OpenReader(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read Excel file"})
		return
	}

	rows, err := xlFile.GetRows("Sheet1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rows"})
		return
	}

	fmt.Println("this are the emails lists", rows)

	var emails []string
	for _, row := range rows {
		if len(row) > 0 {
			emails = append(emails, row[0])
		}
	}

	for i := 0; i < len(emails); i += 100 {
		end := i + 100
		if end > len(emails) {
			end = len(emails)
		}

		batch := emails[i:end]
		for _, email := range batch {
			err := SendEmail(email, "Subscriber", "Newsletter", "templates/newsletter.html", nil)
			if err != nil {
				fmt.Println("Failed to send email to:", email, "Error:", err)
			}
		}

		time.Sleep(1 * time.Minute)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Emails sent successfully"})
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
package auth

import (
	"example/go-api/orm"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var hmacSampleSecret []byte

type RegisterForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
}
type LoginForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type DeleteForm struct {
	Filename string `json:"filename" binding:"required"`
}
type RenameForm struct {
	FilenameOri string `json:"filenameori" binding:"required"`
	FilenameNew string `json:"filenamenew" binding:"required"`
}

func Register(c *gin.Context) {
	var json RegisterForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var userExist orm.User
	orm.DB.Where("username = ?", json.Username).First(&userExist)
	if userExist.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "User Exist",
		})
		return
	}

	user := orm.User{Username: json.Username, Password: json.Password, Fullname: json.Fullname}
	orm.DB.Create(&user)
	if user.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": json,
			"userId":  user.ID,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "error",
		})
	}
}

func Login(c *gin.Context) {
	var json LoginForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var userExist orm.User
	orm.DB.Where("username = ?", json.Username).First(&userExist)
	if userExist.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "User Dos Not Exist",
		})
		return
	}
	if userExist.Password == json.Password {

		hmacSampleSecret = []byte("my_secret_key")
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": userExist.ID,
		})
		tokenString, _ := token.SignedString(hmacSampleSecret)
		c.JSON(http.StatusOK, gin.H{
			"message": "Login Success", "token": tokenString,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Login Failed",
		})
	}
}

func UploadFile(c *gin.Context) {
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}

	filename := filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, fmt.Sprintf("./images/%s", filename)); err != nil {
		c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
		return
	}

	c.String(http.StatusOK, "File %s uploaded successfully with fields", file.Filename)
}
func ListImg(c *gin.Context) {
	tagsList := []string{}
	files, err := ioutil.ReadDir("./images")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		tagsList = append(tagsList, file.Name())
	}
	c.JSON(http.StatusOK, gin.H{
		"Image": tagsList,
	})

}
func Delete(c *gin.Context) {
	var json DeleteForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	e := os.Remove("./images/" + json.Filename)
	if e != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "The system cannot find the file specified",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Delete Success",
		})
	}
}

func Rename(c *gin.Context) {
	var json RenameForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	Original_Path := "./images/" + json.FilenameOri
	New_Path := "./images/" + json.FilenameNew
	e := os.Rename(Original_Path, New_Path)
	if e != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "The system cannot find the file specified",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Rename Success",
		})
	}

}

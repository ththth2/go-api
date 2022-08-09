package main

import (
	Controller "example/go-api/controller/auth"
	"example/go-api/orm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	orm.Connect()

	r := gin.Default()

	r.Use(cors.Default())
	r.GET("/list", Controller.ListImg)
	r.POST("/register", Controller.Register)
	r.POST("/login", Controller.Login)
	r.POST("/api/upload", Controller.UploadFile)
	r.POST("/api/delete", Controller.Delete)
	r.POST("/api/rename", Controller.Rename)
	r.Run("localhost:8080")
}

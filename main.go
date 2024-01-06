package main

import (
	"example/Studying/controllers"
	"example/Studying/initializers"

	_ "example/Studying/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()

	r.POST("/todo", controllers.ToDoCreate)
	r.PUT("/todo/:id", controllers.ToDoUpdate)
	r.DELETE("/todo/:id", controllers.ToDoDelete)

	r.GET("/todo", controllers.ToDoIndex)
	r.GET("/todo/:id", controllers.ToDoShow)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Run()
}

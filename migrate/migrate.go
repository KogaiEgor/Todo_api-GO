package main

import (
	"example/Studying/initializers"
	"example/Studying/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.ToDo{})
}

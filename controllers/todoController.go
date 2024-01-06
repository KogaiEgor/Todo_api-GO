package controllers

import (
	"example/Studying/initializers"
	"example/Studying/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ToDoCreate godoc
// @Summary Create a new todo
// @Description Create a new todo with the input payload
// @Tags todos
// @Accept  json
// @Produce  json
// @Param title body string true "Title of the Todo"
// @Param body body string true "Body of the Todo"
// @Param status body bool true "Status of the Todo"
// @Success 201 {object} models.ToDo "Successfully created"
// @Failure 400 {string} string "Invalid request"
// @Router /todos [post]
func ToDoCreate(c *gin.Context) {
	//Get data
	var body struct {
		Title  string
		Body   string
		Status bool
	}

	c.Bind(&body)

	//Create a ToDo
	todo := models.ToDo{Title: body.Title, Body: body.Body, Status: body.Status}
	result := initializers.DB.Create(&todo)

	if result.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	//Return data
	c.JSON(http.StatusCreated, gin.H{
		"todo": todo,
	})
}

// ToDoIndex godoc
// @Summary List todos
// @Description Get a list of todos, optionally filtered by status
// @Tags todos
// @Accept  json
// @Produce  json
// @Param status query bool false "Filter by status"
// @Success 200 {array} models.ToDo "List of todos"
// @Failure 500 {string} string "Internal server error"
// @Router /todos [get]
func ToDoIndex(c *gin.Context) {
	//Get data
	var todos []models.ToDo
	status := c.Query("status")

	//check param
	if status != "" {
		filteredStatus, err := strconv.ParseBool(status)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong status format"})
			return
		}
		result := initializers.DB.Where("status = ?", filteredStatus).Find(&todos)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
	} else {
		initializers.DB.Find(&todos)
	}

	//Respond with data
	c.JSON(http.StatusOK, gin.H{
		"todos": todos,
	})
}

// ToDoShow godoc
// @Summary Show a todo
// @Description Get details of a todo by ID
// @Tags todos
// @Accept  json
// @Produce  json
// @Param id path int true "Todo ID"
// @Success 200 {object} models.ToDo "Todo details"
// @Failure 404 {string} string "Todo not found"
// @Router /todos/{id} [get]
func ToDoShow(c *gin.Context) {
	//Get param
	id := c.Param("id")

	//Get todo
	var todo models.ToDo
	result := initializers.DB.Find(&todo, id)

	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Элемент не найден"})
		return
	}

	//Respond with single todo
	c.JSON(http.StatusOK, gin.H{
		"todo": todo,
	})
}

// ToDoUpdate godoc
// @Summary Update a todo
// @Description Update a todo with the specified ID
// @Tags todos
// @Accept  json
// @Produce  json
// @Param id path int true "Todo ID"
// @Param title body string false "Title of the Todo"
// @Param body body string false "Body of the Todo"
// @Param status body bool false "Status of the Todo"
// @Success 200 {object} models.ToDo "Successfully updated"
// @Failure 404 {string} string "Todo not found"
// @Router /todos/{id} [put]
func ToDoUpdate(c *gin.Context) {
	//Get param
	id := c.Param("id")

	//Get todo
	var todo models.ToDo
	result := initializers.DB.Find(&todo, id)

	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Элемент не найден"})
		return
	}

	//Get data
	var body struct {
		Title  string
		Body   string
		Status bool
	}

	c.Bind(&body)

	//Update todo
	initializers.DB.Model(&todo).Updates(map[string]interface{}{
		"Title":  body.Title,
		"Body":   body.Body,
		"Status": body.Status,
	})

	//Respond with updated todo
	c.JSON(http.StatusOK, gin.H{
		"todo": todo,
	})
}

// ToDoDelete godoc
// @Summary Delete a todo
// @Description Delete a todo with the specified ID
// @Tags todos
// @Accept  json
// @Produce  json
// @Param id path int true "Todo ID"
// @Success 204 {string} string "Successfully deleted"
// @Failure 404 {string} string "Todo not found"
// @Router /todos/{id} [delete]
func ToDoDelete(c *gin.Context) {
	//Get param
	id := c.Param("id")

	//Delete todo
	result := initializers.DB.Delete(&models.ToDo{}, id)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "ToDo wasn't found"})
		return
	}

	//Respond with status
	c.Status(http.StatusNoContent)
}

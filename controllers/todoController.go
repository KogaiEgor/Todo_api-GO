package controllers

import (
	"example/Studying/initializers"
	"example/Studying/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type body struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Status bool   `json:"status"`
}

// ToDoCreate godoc
// @Summary Create a new todo
// @Description Создание нового todo
// @Description title должен быть не короче 3
// @Tags todos
// @Accept  json
// @Produce  json
// @Param todo body body true "Create Todo"
// @Success 201 {object} models.ToDo "Successfully created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /todo [post]
func ToDoCreate(c *gin.Context) {
	//Get data
	var body body

	c.Bind(&body)

	if len(body.Title) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный формат todo"})
		return
	}

	//Create a ToDo
	todo := models.ToDo{Title: body.Title, Body: body.Body, Status: body.Status}
	result := initializers.DB.Create(&todo)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный формат todo"})
		return
	}

	//Return data
	c.JSON(http.StatusCreated, gin.H{
		"todo": todo,
	})
}

// ToDoIndex godoc
// @Summary List todos
// @Description Получение списка todo, опционально можно отфильтровать по статусу
// @Tags todos
// @Accept  json
// @Produce  json
// @Param status query bool false "Filter by status"
// @Success 200 {array} models.ToDo "List of todos"
// @Failure 500 {string} string "Internal server error"
// @Router /todo [get]
func ToDoIndex(c *gin.Context) {
	//Get data
	var todos []models.ToDo
	status := c.Query("status")

	//check param
	if status != "" {
		filteredStatus, err := strconv.ParseBool(status)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный формат статуса"})
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
// @Description Получение todo по id
// @Tags todos
// @Accept  json
// @Produce  json
// @Param id path int true "Todo ID"
// @Success 200 {object} models.ToDo "Todo details"
// @Failure 404 {object} map[string]string "Todo not found"
// @Router /todo/{id} [get]
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
// @Description Обновление todo по id
// @Tags todos
// @Accept  json
// @Produce  json
// @Param id path int true "Todo ID"
// @Param todo body body false "Update Todo"
// @Success 200 {object} models.ToDo "Successfully updated"
// @Failure 404 {object} map[string]string "Todo not found"
// @Router /todo/{id} [put]
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
	var body body

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
// @Description Удаление todo по id
// @Tags todos
// @Accept  json
// @Produce  json
// @Param id path int true "Todo ID"
// @Success 204 {string} string "Successfully deleted"
// @Failure 404 {object} map[string]string "Todo not found"
// @Router /todo/{id} [delete]
func ToDoDelete(c *gin.Context) {
	//Get param
	id := c.Param("id")

	//Delete todo
	result := initializers.DB.Delete(&models.ToDo{}, id)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Элемент не найден"})
		return
	}

	//Respond with status
	c.Status(http.StatusNoContent)
}

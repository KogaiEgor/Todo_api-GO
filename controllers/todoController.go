package controllers

import (
	"example/Studying/models"
	"example/Studying/services"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title have to be more than 3 letters"})
		return
	}

	// Create a ToDo using the service
	todoService := services.NewTodoService()
	todo, err := todoService.CreateTodo(body.Title, body.Body, body.Status)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong todo format"})
		return
	}

	// Return data
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
	todoService := services.NewTodoService()

	//check param
	if status != "" {
		filteredStatus, err := strconv.ParseBool(status)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong todo format"})
			return
		}

		todos, err = todoService.FindByStatus(filteredStatus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		var err error
		todos, err = todoService.GetAllTodos()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
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

	todoService := services.NewTodoService()

	todo, err := todoService.FindTodo(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ToDo doesn't exist"})
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
	todoService := services.NewTodoService()
	todo, err := todoService.FindTodo(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ToDo doesn't exist"})
		return
	}

	//Get data
	var body body

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong data format"})
		return
	}

	//Update todo
	if err := todoService.UpdateTodo(todo, body.Title, body.Body, body.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
	// Get param
	id := c.Param("id")
	todoService := services.NewTodoService()

	// Delete todo using the service
	if err := todoService.DeleteTodo(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Respond with status
	c.Status(http.StatusNoContent)
}

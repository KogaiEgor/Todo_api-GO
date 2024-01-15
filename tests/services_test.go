package main

import (
	"example/Studying/initializers"
	"example/Studying/models"
	"example/Studying/services"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllTodos(t *testing.T) {
	todoService := services.NewTodoService()

	testTodos := []models.ToDo{
		{Title: "Test ToDo 1", Body: "Test Body 1", Status: true},
		{Title: "Test ToDo 2", Body: "Test Body 2", Status: false},
		{Title: "Test ToDo 3", Body: "Test Body 3", Status: true},
		{Title: "Test ToDo 4", Body: "Test Body 4", Status: false},
		{Title: "Test ToDo 5", Body: "Test Body 5", Status: true},
		{Title: "Test ToDo 6", Body: "Test Body 6", Status: false},
		{Title: "Test ToDo 7", Body: "Test Body 7", Status: true},
		{Title: "Test ToDo 8", Body: "Test Body 8", Status: false},
	}

	todos, err := todoService.GetAllTodos()
	assert.NoError(t, err)
	assert.NotNil(t, todos)
	assert.Len(t, todos, len(testTodos))
	for _, expectedTodo := range testTodos {
		found := false
		for _, actualTodo := range todos {
			if actualTodo.Title == expectedTodo.Title && actualTodo.Body == expectedTodo.Body && actualTodo.Status == expectedTodo.Status {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected todo not found: %+v", expectedTodo)
	}

	fmt.Println("TestGetAllTodos passed!")
}

func TestCreateTodo(t *testing.T) {
	todoService := services.NewTodoService()
	title := "Test Todo"
	body := "This is a test todo"
	status := false

	todo, err := todoService.CreateTodo(title, body, status)

	assert.NoError(t, err)
	assert.NotNil(t, todo)
	assert.Equal(t, title, todo.Title)
	assert.Equal(t, body, todo.Body)
	assert.Equal(t, status, todo.Status)

	fmt.Println("TestCreateTodo passed!")
}

func TestFindTodo(t *testing.T) {
	todoService := services.NewTodoService()

	testTodo := models.ToDo{Title: "Test ToDo 1", Body: "Test Body 1", Status: true}
	foundTodo, err := todoService.FindTodo("1")
	assert.NoError(t, err)
	assert.NotNil(t, foundTodo)
	assert.Equal(t, testTodo.Title, foundTodo.Title)
	assert.Equal(t, testTodo.Body, foundTodo.Body)
	assert.Equal(t, testTodo.Status, foundTodo.Status)

	// Проверка нахождения несуществующей задачи
	_, err = todoService.FindTodo("123123")
	assert.Error(t, err)

	fmt.Println("TestFindTodo passed!")
}

func TestFindByStatus(t *testing.T) {
	todoService := services.NewTodoService()

	todos_true, err := todoService.FindByStatus(true)

	assert.NoError(t, err)
	assert.NotNil(t, todos_true)
	for _, line := range todos_true {
		assert.True(t, line.Status)
	}

	todos_false, err := todoService.FindByStatus(false)

	assert.NoError(t, err)
	assert.NotNil(t, todos_true)
	for _, line := range todos_false {
		assert.False(t, line.Status)
	}

	fmt.Println("TestFindByStatus passed!")
}

func TestUpdateTodo(t *testing.T) {
	todoService := services.NewTodoService()

	testTodo := models.ToDo{Title: "Original Title", Body: "Original Body", Status: false}
	if err := initializers.DB.Create(&testTodo).Error; err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	updatedTitle := "Updated Title"
	updatedBody := "Updated Body"
	updatedStatus := true
	err := todoService.UpdateTodo(&testTodo, updatedTitle, updatedBody, updatedStatus)
	assert.NoError(t, err)

	var updatedTodo models.ToDo
	initializers.DB.First(&updatedTodo, testTodo.ID)
	assert.Equal(t, updatedTitle, updatedTodo.Title)
	assert.Equal(t, updatedBody, updatedTodo.Body)
	assert.Equal(t, updatedStatus, updatedTodo.Status)

	fmt.Println("TestUpdateTodo passed!")
}

func TestDeleteTodo(t *testing.T) {
	todoService := services.NewTodoService()

	testTodo := models.ToDo{Title: "Test Todo", Body: "Test Body", Status: false}
	if err := initializers.DB.Create(&testTodo).Error; err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	err := todoService.DeleteTodo(fmt.Sprintf("%v", testTodo.ID))
	assert.NoError(t, err)

	var deletedTodo models.ToDo
	result := initializers.DB.First(&deletedTodo, testTodo.ID)
	assert.Error(t, result.Error)

	err = todoService.DeleteTodo("99999")
	assert.Error(t, err)

	fmt.Println("TestDeleteTodo passed!")
}

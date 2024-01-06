package main

import (
	"bytes"
	"encoding/json"
	"example/Studying/controllers"
	"example/Studying/initializers"
	"example/Studying/models"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("TEST_DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	fmt.Println("DSN:", dsn)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database")
	}
}

func InitializeTestDB(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.ToDo{}); err != nil {
		return err
	}

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

	for _, todo := range testTodos {
		if result := db.Create(&todo); result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func ClearTestDB(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE to_dos RESTART IDENTITY CASCADE")
}

func TestMain(m *testing.M) {
	if _, err := os.Stat("../.env"); err == nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	ConnectToDB()

	if err := InitializeTestDB(DB); err != nil {
		log.Fatalf("Ошибка при инициализации тестовой базы данных: %s", err)
	}

	initializers.DB = DB

	code := m.Run()

	ClearTestDB(DB)

	os.Exit(code)
}

func TestToDoIndexWithoutStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/todos", func(c *gin.Context) {
		controllers.ToDoIndex(c)
	})

	req, _ := http.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Todos []models.ToDo `json:"todos"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Todos)
}

func TestToDoIndexWithStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/todos", func(c *gin.Context) {
		controllers.ToDoIndex(c)
	})

	reqWithStatus, _ := http.NewRequest("GET", "/todos?status=true", nil)
	wWithStatus := httptest.NewRecorder()
	r.ServeHTTP(wWithStatus, reqWithStatus)

	assert.Equal(t, http.StatusOK, wWithStatus.Code)

	var responseWithStatus struct {
		Todos []models.ToDo `json:"todos"`
	}
	err := json.Unmarshal(wWithStatus.Body.Bytes(), &responseWithStatus)
	assert.NoError(t, err)
	assert.NotEmpty(t, responseWithStatus.Todos)

	for _, todo := range responseWithStatus.Todos {
		assert.True(t, todo.Status, "Все возвращенные ToDo должны иметь статус true")
	}
}

func TestToDoCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/todos", func(c *gin.Context) {
		controllers.ToDoCreate(c)
	})

	newToDo := models.ToDo{
		Title:  "New Test ToDo",
		Body:   "Test ToDo Body",
		Status: true,
	}
	jsonData, _ := json.Marshal(newToDo)

	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		ToDo models.ToDo `json:"todo"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, newToDo.Title, response.ToDo.Title)
	assert.Equal(t, newToDo.Body, response.ToDo.Body)
	assert.Equal(t, newToDo.Status, response.ToDo.Status)
}

func TestToDoCreateWrongData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/todos", func(c *gin.Context) {
		controllers.ToDoCreate(c)
	})

	newToDo := models.ToDo{
		Title:  "N",
		Body:   "Test ToDo Body",
		Status: true,
	}
	jsonData, _ := json.Marshal(newToDo)

	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestToDoShow(t *testing.T) {
	testToDo := models.ToDo{Title: "Test ToDo", Body: "Test Body", Status: true}
	if result := initializers.DB.Create(&testToDo); result.Error != nil {
		t.Fatalf("Ошибка при создании тестовой задачи: %s", result.Error)
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/todos/:id", func(c *gin.Context) {
		controllers.ToDoShow(c)
	})

	req, _ := http.NewRequest("GET", fmt.Sprintf("/todos/%d", testToDo.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		ToDo models.ToDo `json:"todo"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, testToDo.ID, response.ToDo.ID)
	assert.Equal(t, testToDo.Title, response.ToDo.Title)
	assert.Equal(t, testToDo.Body, response.ToDo.Body)
	assert.Equal(t, testToDo.Status, response.ToDo.Status)

	reqNotFound, _ := http.NewRequest("GET", "/todos/99999", nil)
	wNotFound := httptest.NewRecorder()
	r.ServeHTTP(wNotFound, reqNotFound)

	assert.Equal(t, http.StatusNotFound, wNotFound.Code)
}

func TestToDoUpdate(t *testing.T) {
	testToDo := models.ToDo{Title: "Original Title", Body: "Original Body", Status: true}
	if result := initializers.DB.Create(&testToDo); result.Error != nil {
		t.Fatalf("Ошибка при создании тестовой задачи: %s", result.Error)
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.PUT("/todos/:id", func(c *gin.Context) {
		controllers.ToDoUpdate(c)
	})

	updatedToDo := models.ToDo{
		Title:  "Updated Title",
		Body:   "Updated Body",
		Status: false,
	}
	jsonData, _ := json.Marshal(updatedToDo)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/todos/%d", testToDo.ID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		ToDo models.ToDo `json:"todo"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, updatedToDo.Title, response.ToDo.Title)
	assert.Equal(t, updatedToDo.Body, response.ToDo.Body)
	assert.Equal(t, updatedToDo.Status, response.ToDo.Status)

	reqNotFound, _ := http.NewRequest("PUT", "/todos/99999", bytes.NewBuffer(jsonData))
	reqNotFound.Header.Set("Content-Type", "application/json")
	wNotFound := httptest.NewRecorder()
	r.ServeHTTP(wNotFound, reqNotFound)

	assert.Equal(t, http.StatusNotFound, wNotFound.Code)
}

func TestToDoDelete(t *testing.T) {
	testToDo := models.ToDo{Title: "Test ToDo", Body: "Test Body", Status: true}
	if result := initializers.DB.Create(&testToDo); result.Error != nil {
		t.Fatalf("Ошибка при создании тестовой задачи: %s", result.Error)
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.DELETE("/todos/:id", func(c *gin.Context) {
		controllers.ToDoDelete(c)
	})

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/todos/%d", testToDo.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	reqNotFound, _ := http.NewRequest("DELETE", "/todos/99999", nil)
	wNotFound := httptest.NewRecorder()
	r.ServeHTTP(wNotFound, reqNotFound)

	assert.Equal(t, http.StatusNotFound, wNotFound.Code)
}

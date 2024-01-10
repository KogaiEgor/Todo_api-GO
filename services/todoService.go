package services

import (
	"example/Studying/initializers"
	"example/Studying/models"
	"fmt"
)

type TodoService struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Status bool   `json:"status"`
}

func NewTodoService() *TodoService {
	return &TodoService{}
}

func (s *TodoService) CreateTodo(title string, body string, status bool) (*models.ToDo, error) {
	todo := &models.ToDo{
		Title:  title,
		Body:   body,
		Status: status,
	}

	if err := initializers.DB.Create(todo).Error; err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *TodoService) GetAllTodos() ([]models.ToDo, error) {
	var todos []models.ToDo

	if err := initializers.DB.Find(&todos).Error; err != nil {
		return nil, err
	}

	return todos, nil
}

func (s *TodoService) FindTodo(todoID string) (*models.ToDo, error) {
	var todo models.ToDo

	if err := initializers.DB.First(&todo, todoID).Error; err != nil {
		return nil, err
	}

	return &todo, nil
}

func (s *TodoService) FindByStatus(status bool) ([]models.ToDo, error) {
	var todos []models.ToDo

	if err := initializers.DB.Where("status = ?", status).Find(&todos).Error; err != nil {
		return nil, err
	}

	return todos, nil
}

func (s *TodoService) UpdateTodo(todo *models.ToDo, title string, body string, status bool) error {
	updateFields := map[string]interface{}{
		"Title":  title,
		"Body":   body,
		"Status": status,
	}

	if err := initializers.DB.Model(todo).Updates(updateFields).Error; err != nil {
		return err
	}

	return nil
}

func (s *TodoService) DeleteTodo(todoID string) error {
	result := initializers.DB.Delete(&models.ToDo{}, todoID)

	if result.RowsAffected == 0 {
		return fmt.Errorf("There is no todo with id %s", todoID)
	}

	return nil
}

package services

import (
	"boilerplate-api/api/repository"
	"boilerplate-api/models"
	"boilerplate-api/utils"
)

// TodoService -> struct
type TodoService struct {
	repository repository.TodoRepository
}

// NewTodoService  -> creates a new Todoservice
func NewTodoService(repository repository.TodoRepository) TodoService {
	return TodoService{
		repository: repository,
	}
}

// CreateTodo -> call to create the Todo
func (c TodoService) CreateTodo(todo models.Todo) (models.Todo, error) {
	return c.repository.Create(todo)
}

// GetAllTodo -> call to create the Todo
func (c TodoService) GetAllTodo(pagination utils.Pagination) ([]models.Todo, int64, error) {
	return c.repository.GetAllTodo(pagination)
}

// GetOneTodo -> Get One Todo By Id
func (c TodoService) GetOneTodo(ID int64) (models.Todo, error) {
	return c.repository.GetOneTodo(ID)
}

// UpdateOneTodo -> Update One Todo By Id
func (c TodoService) UpdateOneTodo(todo models.Todo) error {
	return c.repository.UpdateOneTodo(todo)
}

// DeleteOneTodo -> Delete One Todo By Id
func (c TodoService) DeleteOneTodo(ID int64) error {
	return c.repository.DeleteOneTodo(ID)

}

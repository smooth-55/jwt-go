package repository

import (
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
)

// TodoRepository database structure
type TodoRepository struct {
	db     infrastructure.Database
	logger infrastructure.Logger
}

// NewTodoRepository creates a new Todo repository
func NewTodoRepository(db infrastructure.Database, logger infrastructure.Logger) TodoRepository {
	return TodoRepository{
		db:     db,
		logger: logger,
	}
}

// Create Todo
func (c TodoRepository) Create(Todo models.Todo) (models.Todo, error) {
	return Todo, c.db.DB.Create(&Todo).Error
}

// GetAllTodo -> Get All todos
func (c TodoRepository) GetAllTodo(pagination utils.Pagination) ([]models.Todo, int64, error) {
	var todos []models.Todo
	var totalRows int64 = 0
	queryBuider := c.db.DB.Model(&models.Todo{}).Offset(pagination.Offset).Order(pagination.Sort)

	if !pagination.All {
		queryBuider = queryBuider.Limit(pagination.PageSize)
	}

	err := queryBuider.
		Find(&todos).
		Offset(-1).
		Limit(-1).
		Count(&totalRows).Error
	return todos, totalRows, err
}

// GetOneTodo -> Get One Todo By Id
func (c TodoRepository) GetOneTodo(ID int64) (models.Todo, error) {
	Todo := models.Todo{}
	return Todo, c.db.DB.
		Where("id = ?", ID).First(&Todo).Error
}

// UpdateOneTodo -> Update One Todo By Id
func (c TodoRepository) UpdateOneTodo(Todo models.Todo) error {
	return c.db.DB.Model(&models.Todo{}).
		Where("id = ?", Todo.ID).
		Updates(Todo).Error
}

// DeleteOneTodo -> Delete One Todo By Id
func (c TodoRepository) DeleteOneTodo(ID int64) error {
	return c.db.DB.
		Where("id = ?", ID).
		Delete(&models.Todo{}).
		Error
}

package controllers

import (
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TodoController -> struct
type TodoController struct {
	logger      infrastructure.Logger
	TodoService services.TodoService
}

// NewTodoController -> constructor
func NewTodoController(
	logger infrastructure.Logger,
	TodoService services.TodoService,
) TodoController {
	return TodoController{
		logger:      logger,
		TodoService: TodoService,
	}
}

// CreateTodo -> Create Todo
func (cc TodoController) CreateTodo(c *gin.Context) {
	todo := models.Todo{}

	if err := c.ShouldBindJSON(&todo); err != nil {
		cc.logger.Zap.Error("Error [CreateTodo] (ShouldBindJson) : ", err)
		err := errors.BadRequest.Wrap(err, "Failed to bind Todo")
		responses.HandleError(c, err)
		return
	}

	if _, err := cc.TodoService.CreateTodo(todo); err != nil {
		cc.logger.Zap.Error("Error [CreateTodo] [db CreateTodo]: ", err.Error())
		err := errors.BadRequest.Wrap(err, "Failed To Create Todo")
		responses.HandleError(c, err)
		return
	}

	responses.SuccessJSON(c, http.StatusOK, "Todo Created Sucessfully")
}

// GetAllTodo -> Get All Todo
func (cc TodoController) GetAllTodo(c *gin.Context) {

	pagination := utils.BuildPagination(c)
	pagination.Sort = "CreateDateTime desc"
	todos, count, err := cc.TodoService.GetAllTodo(pagination)

	if err != nil {
		cc.logger.Zap.Error("Error finding Todo records", err.Error())
		err := errors.InternalError.Wrap(err, "Failed To Find Todo")
		responses.HandleError(c, err)
		return
	}
	responses.JSONCount(c, http.StatusOK, todos, count)

}

// GetOneTodo -> Get One Todo
func (cc TodoController) GetOneTodo(c *gin.Context) {
	ID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	todo, err := cc.TodoService.GetOneTodo(ID)

	if err != nil {
		cc.logger.Zap.Error("Error [GetOneTodo] [db GetOneTodo]: ", err.Error())
		err := errors.InternalError.Wrap(err, "Failed To Find Todo")
		responses.HandleError(c, err)
		return
	}
	responses.JSON(c, http.StatusOK, todo)

}

// UpdateOneTodo -> Update One Todo By Id
func (cc TodoController) UpdateOneTodo(c *gin.Context) {
	ID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	todo := models.Todo{}

	if err := c.ShouldBindJSON(&todo); err != nil {
		cc.logger.Zap.Error("Error [UpdateTodo] (ShouldBindJson) : ", err)
		err := errors.BadRequest.Wrap(err, "failed to update todo")
		responses.HandleError(c, err)
		return
	}
	todo.ID = ID

	if err := cc.TodoService.UpdateOneTodo(todo); err != nil {
		cc.logger.Zap.Error("Error [UpdateTodo] [db UpdateTodo]: ", err.Error())
		err := errors.InternalError.Wrap(err, "failed to update todo")
		responses.HandleError(c, err)
		return
	}

	responses.SuccessJSON(c, http.StatusOK, "Todo Updated Sucessfully")
}

// DeleteOneTodo -> Delete One Todo By Id
func (cc TodoController) DeleteOneTodo(c *gin.Context) {
	ID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	err := cc.TodoService.DeleteOneTodo(ID)

	if err != nil {
		cc.logger.Zap.Error("Error [DeleteOneTodo] [db DeleteOneTodo]: ", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to Delete Todo")
		responses.HandleError(c, err)
		return
	}

	responses.SuccessJSON(c, http.StatusOK, "Todo Deleted Sucessfully")
}

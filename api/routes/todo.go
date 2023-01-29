package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/api/middlewares"
	"boilerplate-api/infrastructure"
)

// TodoRoutes -> struct
type TodoRoutes struct {
	logger         infrastructure.Logger
	router         infrastructure.Router
	todoController controllers.TodoController
	middleware     middlewares.FirebaseAuthMiddleware
}

// NewTodoRoutes -> creates new Todo controller
func NewTodoRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	todoController controllers.TodoController,
	middleware middlewares.FirebaseAuthMiddleware,
) TodoRoutes {
	return TodoRoutes{
		router:         router,
		logger:         logger,
		todoController: todoController,
		middleware:     middleware,
	}
}

// Setup todo routes
func (c TodoRoutes) Setup() {
	c.logger.Zap.Info(" Setting up Todo routes")
	todo := c.router.Gin.Group("/todo")
	{
		todo.POST("", c.todoController.CreateTodo)
		todo.GET("", c.todoController.GetAllTodo)
		todo.GET("/:id", c.todoController.GetOneTodo)
		todo.PUT("/:id", c.todoController.UpdateOneTodo)
		todo.DELETE("/:id", c.todoController.DeleteOneTodo)
	}
}

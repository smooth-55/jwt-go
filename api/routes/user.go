package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/api/middlewares"
	"boilerplate-api/infrastructure"
)

// UserRoutes -> struct
type UserRoutes struct {
	logger            infrastructure.Logger
	router            infrastructure.Router
	userController    controllers.UserController
	middleware        middlewares.FirebaseAuthMiddleware
	trxMiddleware     middlewares.DBTransactionMiddleware
	jwtAuthMiddleware middlewares.JWTAuthMiddleWare
}

// Setup user routes
func (i UserRoutes) Setup() {
	i.logger.Zap.Info(" Setting up user routes")
	users := i.router.Gin.Group("/user").Use(i.jwtAuthMiddleware.Handle())
	{
		users.GET("", i.userController.GetAllUsers)
		users.GET("/:id", i.userController.GetOneUser)
		users.PUT("/:id", i.trxMiddleware.DBTransactionHandle(), i.userController.UpdateUser)
		users.DELETE("/:id", i.trxMiddleware.DBTransactionHandle(), i.userController.DeleteOneUser)
		users.POST("", i.trxMiddleware.DBTransactionHandle(), i.userController.CreateUser)
	}
	user := i.router.Gin.Group("/jwt-login")
	{
		user.POST("", i.userController.LoginUser)

	}

}

// NewUserRoutes -> creates new user controller
func NewUserRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	userController controllers.UserController,
	middleware middlewares.FirebaseAuthMiddleware,
	trxMiddleware middlewares.DBTransactionMiddleware,
	jwtAuthMiddleware middlewares.JWTAuthMiddleWare,
) UserRoutes {
	return UserRoutes{
		router:            router,
		logger:            logger,
		userController:    userController,
		middleware:        middleware,
		trxMiddleware:     trxMiddleware,
		jwtAuthMiddleware: jwtAuthMiddleware,
	}
}

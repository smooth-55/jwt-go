package controllers

import (
	"boilerplate-api/api/middlewares"
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/api/validators"
	"boilerplate-api/constants"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserController -> struct
type UserController struct {
	logger          infrastructure.Logger
	userService     services.UserService
	env             infrastructure.Env
	validator       validators.UserValidator
	firebaseService services.FirebaseService
}

// NewUserController -> constructor
func NewUserController(
	logger infrastructure.Logger,
	userService services.UserService,
	env infrastructure.Env,
	validator validators.UserValidator,
	firebaseService services.FirebaseService,
) UserController {
	return UserController{
		logger:          logger,
		userService:     userService,
		env:             env,
		validator:       validator,
		firebaseService: firebaseService,
	}
}

// CreateUser -> Create User
func (cc UserController) CreateUser(c *gin.Context) {
	reqData := struct {
		models.User
		ConfirmPassword string `json:"confirm_password" validate:"required"`
	}{}
	trx := c.MustGet(constants.DBTransaction).(*gorm.DB)
	if err := c.ShouldBindJSON(&reqData); err != nil {
		cc.logger.Zap.Error("Error [CreateUser] (ShouldBindJson) : ", err)
		err := errors.BadRequest.Wrap(err, "Failed to bind user data")
		responses.HandleError(c, err)
		return
	}
	if reqData.User.Password != reqData.ConfirmPassword {
		cc.logger.Zap.Error("Password and confirm password not matching : ")
		responses.ErrorJSON(c, http.StatusBadRequest, "Password and confirm password should be same.")
		return
	}

	if validationErr := cc.validator.Validate.Struct(reqData); validationErr != nil {
		err := errors.BadRequest.Wrap(validationErr, "Validation error")
		err = errors.SetCustomMessage(err, "Invalid input information")
		err = errors.AddErrorContextBlock(err, cc.validator.GenerateValidationResponse(validationErr))
		responses.HandleError(c, err)
		return
	}

	if !utils.IsValidEmail(reqData.User.Email) {
		cc.logger.Zap.Error("Invalid email")
		responses.ErrorJSON(c, http.StatusBadRequest, "Invalid Email")
		return
	}
	password, err := bcrypt.GenerateFromPassword([]byte(reqData.ConfirmPassword), 10)
	reqData.User.Password = string(password)
	if err != nil {
		responses.ErrorJSON(c, http.StatusInternalServerError, "Failed top create hash pw")
		return
	}
	if _, err := cc.userService.WithTrx(trx).CreateUser(&reqData.User); err != nil {
		cc.logger.Zap.Error("Error [CreateUser] [db CreateUser]: ", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to create user")
		responses.HandleError(c, err)
		return
	}

	responses.SuccessJSON(c, http.StatusOK, "user created.")
}

// GetAllUser -> Get All User
func (cc UserController) GetAllUsers(c *gin.Context) {
	pagination := utils.BuildPagination(c)
	users, count, err := cc.userService.GetAllUsers(pagination)
	if err != nil {
		cc.logger.Zap.Error("Error finding user records", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to get users data")
		responses.HandleError(c, err)
		return
	}
	responses.JSONCount(c, http.StatusOK, users, count)
}

func (cc UserController) GetOneUser(c *gin.Context) {
	id := c.Param("id")
	user, err := cc.userService.GetOneUser(id)
	if err != nil {
		cc.logger.Zap.Error("Error finding user records", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to get users data")
		responses.HandleError(c, err)
		return
	}
	responses.SuccessJSON(c, http.StatusOK, user.ToMap())
	return
}

func (cc UserController) DeleteOneUser(c *gin.Context) {
	trx := c.MustGet(constants.DBTransaction).(*gorm.DB)
	f_uid, err := cc.userService.WithTrx(trx).DeleteOneUser(c.Param("id"))
	if err != nil {
		cc.logger.Zap.Error("Error Deleting user record", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to delete user data")
		responses.HandleError(c, err)
		return
	}
	if err := cc.firebaseService.DeleteUser(*f_uid); err != nil {
		cc.logger.Zap.Error("Error Deleting user record from firebase", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to delete user data from firebase")
		responses.HandleError(c, err)
		return
	}
	responses.SuccessJSON(c, http.StatusOK, "User deleted successfully")
	return
}

func (cc UserController) UpdateUser(c *gin.Context) {
	trx := c.MustGet(constants.DBTransaction).(*gorm.DB)
	bodyData := struct {
		Username string `json:"username" validate:"required"`
		FullName string `json:"full_name" validate:"required"`
		Email    string `json:"email" validate:"required"`
		Phone    string `json:"phone" validate:"required"`
		Address  string `json:"address" validate:"required"`
	}{}

	if err := c.ShouldBindJSON(&bodyData); err != nil {
		cc.logger.Zap.Error("Error finding user records", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to get users data")
		responses.HandleError(c, err)
		return
	}
	if validatedError := cc.validator.Validate.Struct(&bodyData); validatedError != nil {
		err := errors.BadRequest.Wrap(validatedError, "Validation error")
		err = errors.SetCustomMessage(err, "Invalid input information")
		err = errors.AddErrorContextBlock(err, cc.validator.GenerateValidationResponse(validatedError))
		responses.HandleError(c, err)
		return
	}
	if !utils.IsValidEmail(bodyData.Email) {
		cc.logger.Zap.Error("Invalid email")
		responses.ErrorJSON(c, http.StatusBadRequest, "Invalid Email")
		return
	}
	fb_user := cc.firebaseService.GetUserByEmail(bodyData.Email)
	if fb_user != "" {
		err := errors.BadRequest.New("Firebase user already exists")
		err = errors.SetCustomMessage(err, "Email address already taken")
		responses.HandleError(c, err)
		return
	}
	bodyDataMap := map[string]interface{}{
		"email":     bodyData.Email,
		"username":  bodyData.Username,
		"phone":     bodyData.Phone,
		"full_name": bodyData.FullName,
		"address":   bodyData.Address,
	}
	user, err := cc.userService.WithTrx(trx).UpdateUser(c.Param("id"), bodyDataMap)
	if err != nil {
		cc.logger.Zap.Error("Error [UpdateUser] [db UpdateUser]: ", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to update user")
		responses.HandleError(c, err)
		return
	}
	userToUpdate := models.UserToUpdate{}
	userToUpdate.Email = user.Email
	userToUpdate.FullName = user.FullName
	userToUpdate.Username = user.Username
	userToUpdate.Address = user.Address
	userToUpdate.Phone = user.Phone
	updateFirebaseUser, err := cc.firebaseService.UpdateUser(user.FirebaseUID, userToUpdate)
	if err != nil {
		cc.logger.Zap.Error("Error [UpdateUser] [db UpdateUser]: ", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to update user")
		responses.HandleError(c, err)
		return
	}
	fmt.Println(updateFirebaseUser)
	responses.SuccessJSON(c, http.StatusOK, user)
	return

}

func (cc UserController) LoginUser(c *gin.Context) {
	var reqData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	// Bind the request payload to a reqData struct
	if err := c.ShouldBindJSON(&reqData); err != nil {
		cc.logger.Zap.Error("Error [ShouldBindJSON] : ", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to bind request data")
		responses.HandleError(c, err)
		return
	}
	user, err := cc.userService.GetOneUserWithEmail(reqData.Email)
	if err != nil {
		responses.ErrorJSON(c, http.StatusBadRequest, "Invalid user credentials1")
		return
	}

	// Check if the password is correct
	// I've encrypted and saved password in DB, so here i am comparing plain text with it's hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqData.Password)); err != nil {
		responses.ErrorJSON(c, http.StatusBadRequest, "Invalid user credentials2")
		return
	}
	userId := fmt.Sprintf("%v", user.ID)

	// Create a new JWT claims object
	claims := middlewares.JWTClaims{
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Subject:   userId,
		},
	}

	// Create a new JWT token using the claims and the secret key
	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, tokenErr := tokenClaim.SignedString([]byte(cc.env.JWT_SECRET))
	if tokenErr != nil {
		responses.ErrorJSON(c, http.StatusInternalServerError, tokenErr.Error())
		return
	}
	data := map[string]interface{}{
		"user":  user.ToMap(),
		"token": token,
	}
	responses.SuccessJSON(c, http.StatusOK, data)
	return
}

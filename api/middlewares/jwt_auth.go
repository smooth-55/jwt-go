package middlewares

import (
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JWTAuthMiddleWare struct {
	jwtService  services.JWTAuthService
	logger      infrastructure.Logger
	env         infrastructure.Env
	db          infrastructure.Database
	userService services.UserService
}

func NewJWTAuthMiddleWare(
	jwtService services.JWTAuthService,
	logger infrastructure.Logger,
	env infrastructure.Env,
	db infrastructure.Database,
	userService services.UserService,

) JWTAuthMiddleWare {
	return JWTAuthMiddleWare{
		jwtService:  jwtService,
		logger:      logger,
		env:         env,
		db:          db,
		userService: userService,
	}
}

type JWTClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func (m JWTAuthMiddleWare) ParseToken(tokenString string) (*jwt.Token, error) {
	// Parse the token using the secret key
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.env.JWT_SECRET), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil

}

func (m JWTAuthMiddleWare) verifyToken(c *gin.Context) (bool, error) {
	// Get the token from the request header
	header := c.GetHeader("Authorization")
	tokenString := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))
	token, err := m.ParseToken(tokenString)
	if err != nil {
		m.logger.Zap.Error("Error parsing token", err.Error())
		return false, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if ok && token.Valid {
		// check if the token expires or not
		if float64(time.Now().Unix()) > float64(claims.ExpiresAt) {
			m.logger.Zap.Error("token expired [ParseWithClaims]")
			err := errors.BadRequest.New("Token already expired")
			err = errors.SetCustomMessage(err, "Token expired")
			return false, err
		}
	} else {
		err := errors.BadRequest.New("Invalid token")
		err = errors.SetCustomMessage(err, "Invalid token")
		return false, err
	}
	// Get user from claims and set
	user, err := m.userService.GetOneUser(claims.Id)
	if err != nil {
		m.logger.Zap.Error("Error finding user records", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to get users data")
		m.logger.Zap.Error("error finding user record")
		return false, err
	}
	// Can set anything in the request context and passes the request to the next handler.
	c.Set("user_id", user.ID)
	return true, nil

}

func (m JWTAuthMiddleWare) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := m.verifyToken(c)
		if !ok {
			m.logger.Zap.Error("Error verifying auth token")
			err = errors.Unauthorized.Wrap(err, "Error verifying auth token")
			err = errors.SetCustomMessage(err, "Unauthorised")
			responses.HandleError(c, err)
			c.Abort()
			return
		}
		c.Next()
	}
}

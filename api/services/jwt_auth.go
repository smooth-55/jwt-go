package services

import (
	"boilerplate-api/infrastructure"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTAuthService struct {
	logger infrastructure.Logger
	env    infrastructure.Env
}

func NewJWTAuthService(
	logger infrastructure.Logger,
	env infrastructure.Env,
) JWTAuthService {
	return JWTAuthService{
		logger: logger,
		env:    env,
	}
}

func (m JWTAuthService) VerifyToken(tokenString string) (*jwt.MapClaims, bool) {

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.env.JWT_SECRET), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return nil, false
		}
	}
	return &claims, true

}

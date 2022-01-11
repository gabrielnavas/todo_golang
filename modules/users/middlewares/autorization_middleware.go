package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type AuthorizationMiddleware interface {
	Autorizate() gin.HandlerFunc
}

type AuthorizationMiddlewareGin struct{}

func NewAuthorizationMiddleware() AuthorizationMiddleware {
	return &AuthorizationMiddlewareGin{}
}

func (middleware *AuthorizationMiddlewareGin) Autorizate() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userId", "12345")
		valueUserId, _ := c.Get("userId")
		userId := valueUserId.(int)

		fmt.Println(userId)
		c.Next()
	}
}

package middlewares

import (
	"api/modules/users/usecases"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	UserId      = "userId"
	LevelAccess = "levelAccess"
)

type AuthorizationMiddleware interface {
	Authorize() gin.HandlerFunc
}

type AuthorizationMiddlewareGin struct {
	tokenManager usecases.TokenManager
}

func NewAuthorizationMiddleware(
	tokenManager usecases.TokenManager,
) AuthorizationMiddleware {
	return &AuthorizationMiddlewareGin{tokenManager}
}

func (middleware *AuthorizationMiddlewareGin) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// split string header Authorization: Authorization: Bearer <token>
		authorizationHeader := c.Request.Header.Get("Authorization")
		if authorizationHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "you need to set header: Authorization: Bearer <token>"})
			return
		}
		authorizationSplited := strings.Split(authorizationHeader, " ")
		if len(authorizationSplited) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "you need to set header: Authorization: Bearer <token>"})
			return
		}

		// get token
		token := authorizationSplited[1]

		// verify token
		payloadToken, err := middleware.tokenManager.VerifyToken(token)
		if err != nil {
			// TODO: make log
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "you don't have authoriation"})
			return
		}

		// add userid and level access from payload token
		c.Set(UserId, payloadToken.UserId)
		c.Set(LevelAccess, payloadToken.LevelAccess)
		c.Next()
	}
}

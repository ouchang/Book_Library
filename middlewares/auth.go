package middlewares

import (
	"gobooklibrary/models"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// https://github.com/appleboy/gin-jwt

const (
	Admin     string = "admin"
	Librarian string = "librarian"
	User      string = "user"
)

func GinJwtMiddlewareHandler(kind string) *jwt.GinJWTMiddleware {
	var identityKey = "id"
	var key []byte

	switch kind {
	case Admin:
		key = []byte("Admin Key")
	case Librarian:
		key = []byte("Librarian Key")
	case User:
		key = []byte("User Key")
	}

	return &jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         key,
		Timeout:     time.Minute * 15,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.LoginAuth); ok {
				return jwt.MapClaims{
					identityKey: v.Login,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.LoginAuth{
				Login: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginAuth models.LoginAuth
			var err error

			if err := c.ShouldBind(&loginAuth); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			// Admin Auth
			if kind == Admin {
				err = models.LoginEmployee(loginAuth, Admin)
			}

			// Librarian
			if kind == Librarian {
				err = models.LoginEmployee(loginAuth, Librarian)
			}

			// User Auth
			if kind == User {
				err = models.LoginUser(loginAuth)
			}

			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return &models.LoginAuth{
				Login: loginAuth.Login,
			}, nil

		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
}

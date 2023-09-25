package routers

import (

	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/SmartfrenDev/go-boilerplate/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

)

func TokenJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			errInfo           []utils.Errors
			tokenAccessString string
		)

		claims := jwt.MapClaims{}
		tokenAccess := c.Request.Header["Authorization"]

		if len(tokenAccess) == 0 {
			errInfo = utils.ErrorWrapper(errInfo, fmt.Sprintf("%v", http.StatusUnauthorized), "token required")
			utils.ResponseWrapperWithErrorInfo(c, nil, nil, http.StatusUnauthorized)
			c.Abort()
			return
		}

		_, err := jwt.ParseWithClaims(tokenAccessString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid token")
			}
			return []byte(os.Getenv("APP_SECRET")), nil
		})

		if err != nil {
			errInfo = utils.ErrorWrapper(errInfo, fmt.Sprintf("%v", http.StatusUnauthorized), err.Error())
			utils.ResponseWrapperWithErrorInfo(c, nil, nil, http.StatusUnauthorized)
			utils.ResponseWrapperWithErrorInfo(c, nil, nil, http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Next()
	}
}

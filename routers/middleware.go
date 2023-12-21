package routers

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/constants"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/personalaccounts"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

func tokenSignature() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			errInfo     []errorsinfo.Errors
			tokenAccess []string
			splitToken  []string
			ok          bool
		)

		claims := jwt.MapClaims{}
		tokenAccess = c.Request.Header["Authorization"]

		if len(tokenAccess) == 0 {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token required")
			response.SendBack(c, struct{}{}, errInfo, http.StatusUnauthorized)
			c.Abort()
			return
		}

		splitToken = strings.Split(tokenAccess[0], "Bearer ")
		token, err := jwt.ParseWithClaims(splitToken[1], claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid token")
			}
			return []byte(os.Getenv("APP_SECRET")), nil
		})

		if err != nil {
			logrus.Error(err.Error())
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "unauthorized token")
			response.SendBack(c, struct{}{}, errInfo, http.StatusUnauthorized)
			c.Abort()
			return
		}

		claims, ok = token.Claims.(jwt.MapClaims)
		if !ok {
			logrus.Error("couldn't parse claims")
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "couldn't parse claims")
			response.SendBack(c, struct{}{}, errInfo, http.StatusUnauthorized)
			c.Abort()
			return
		}

		exp := claims["exp"].(float64)
		if int64(exp) < time.Now().Local().Unix() {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token expired")
			response.SendBack(c, struct{}{}, errInfo, http.StatusUnauthorized)
			c.Abort()
			return
		}

		email := claims["email"].(string)
		c.Set("email", email)
		c.Next()
	}
}

func accountType() gin.HandlerFunc {
	return func(c *gin.Context) {
		var errInfo []errorsinfo.Errors

		usrEmail := c.MustGet("email").(string)
		personalAccount := personalaccounts.Informations(c, usrEmail)

		if personalAccount.ID == uuid.Nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", constants.TokenInvalidInformation)
			response.SendBack(c, struct{}{}, errInfo, http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set("accountType", personalAccount.AccountTypes)
		c.Set("accountID", personalAccount.ID)
		c.Next()
	}
}
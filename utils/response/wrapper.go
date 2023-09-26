package response

import "github.com/gin-gonic/gin"

func SendBack(c *gin.Context, data interface{}, error interface{}, httpStatus int) {
	c.JSONP(httpStatus, gin.H{
		"data":   data,
		"errors": error,
	})
	c.Header("Content-Type", "application/json; charset=utf-8")
}

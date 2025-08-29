package utils

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func ResponseError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"status": "error",
		"message": message,
	})
}



func ResponseSuccess(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{
		"status": "success",
		"data": data,
	})
}

func Logger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lmsgprefix)
}


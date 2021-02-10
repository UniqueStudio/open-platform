package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RecoveryHandler is a func to respond while internal server error occur
func RecoveryHandler(c *gin.Context, err interface{}) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code": http.StatusInternalServerError,
		"err":  err,
	})
}

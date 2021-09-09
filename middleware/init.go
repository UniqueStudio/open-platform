package middleware

import "github.com/gin-gonic/gin"

func InitGlobalMiddleWare(r *gin.Engine) {
	r.Use(tracingMiddleware())
}

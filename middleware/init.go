package middleware

import (
	"github.com/UniqueStudio/open-platform/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func InitGlobalMiddleWare(r *gin.Engine) {
	r.Use(sessions.Sessions("open-platform", utils.RedisSessionStore))
	r.Use(tracingMiddleware())
}

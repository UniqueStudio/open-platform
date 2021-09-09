package middleware

import (
	"errors"
	"net/http"

	"github.com/UniqueStudio/open-platform/config"
	"github.com/UniqueStudio/open-platform/pkg"
	"github.com/UniqueStudio/open-platform/utils"
	"github.com/gin-gonic/gin"
	"github.com/xylonx/zapx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"github.com/gin-contrib/sessions"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apmCtx, span := utils.Tracer.Start(ctx.Request.Context(), "Auth")
		span.SetAttributes(attribute.String("mode", config.Config.Application.Mode))
		// skip when debugging
		if config.Config.Application.Mode == pkg.ModeDebug {
			ctx.Next()
			span.End()
			return
		}

		// judge auth case

		// UI method
		sess := sessions.Default(ctx)
		if au, ok := sess.Get(pkg.SESSION_KEY).(*pkg.AccessUser); ok {
			span.SetAttributes(attribute.Any("accessUser", au))
			span.AddEvent("session auth")
			zapx.WithContext(apmCtx).Info("get accessUser from session")

			ctx.Request = ctx.Request.WithContext(pkg.AddAccessUserIntoContext(apmCtx, au))
			ctx.Next()
			span.End()
			return
		}

		// API method
		accessKey := ctx.GetHeader("AccessKey")
		if accessKey != "" {
			span.SetAttributes(attribute.Any("accessKey", accessKey))
			span.AddEvent("AccessKey auth")
			zapx.WithContext(apmCtx).Info("get accessKey")

			au, err := utils.LoadAccessKey(accessKey)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				zapx.WithContext(apmCtx).Error("load accessUser from accessKey failed", zap.Error(err))

				ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse(err))
				ctx.Abort()
				span.End()
				return
			}

			span.SetAttributes(attribute.Any("accessUser", au))
			zapx.WithContext(apmCtx).Info("get accessUser from accessKey successfully")

			ctx.Request = ctx.Request.WithContext(pkg.AddAccessUserIntoContext(apmCtx, au))
			ctx.Next()
			span.End()
			return
		}

		span.AddEvent("not authorized")
		zapx.WithContext(apmCtx).Warn("not authorized")

		// not auth
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse(errors.New("not authorized")))
		ctx.Abort()
		span.End()
	}
}

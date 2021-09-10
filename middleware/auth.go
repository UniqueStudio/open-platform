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

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apmCtx, span := utils.Tracer.Start(ctx.Request.Context(), "Authentication")
		span.SetAttributes(attribute.String("mode", config.Config.Application.Mode))
		// skip when debugging
		if config.Config.Application.Mode == pkg.ModeDebug {
			// Debug user
			au := &pkg.AccessUser{
				UserID:  "xylonx",
				IsAdmin: true,
			}
			ctx.Request = ctx.Request.WithContext(pkg.AddAccessUserIntoContext(ctx.Request.Context(), au))
			span.End()
			ctx.Next()
			return
		}

		// judge auth case

		// UI method
		sess := sessions.Default(ctx)
		if au, ok := sess.Get(pkg.SESSION_KEY).(*pkg.AccessUser); ok {
			span.SetAttributes(attribute.Any("accessUser", au))
			span.AddEvent("session auth")
			zapx.WithContext(apmCtx).Info("get accessUser from session")

			ctx.Request = ctx.Request.WithContext(pkg.AddAccessUserIntoContext(ctx.Request.Context(), au))
			span.End()
			ctx.Next()
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
				span.End()
				ctx.Abort()
				return
			}

			span.SetAttributes(attribute.Any("accessUser", au))
			zapx.WithContext(apmCtx).Info("get accessUser from accessKey successfully")

			ctx.Request = ctx.Request.WithContext(pkg.AddAccessUserIntoContext(ctx.Request.Context(), au))
			span.End()
			ctx.Next()
			return
		}

		span.AddEvent("not authorized")
		zapx.WithContext(apmCtx).Warn("not authorized")

		// not auth
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse(errors.New("not authorized")))
		span.End()
		ctx.Abort()
	}
}

// FIXME: run order is error
func Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apmCtx, span := utils.Tracer.Start(ctx.Request.Context(), "Authorization")
		au, err := pkg.GetAccessUserFromContext(ctx.Request.Context())
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			zapx.WithContext(apmCtx).Error("get access user from context failed")

			ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse(err))
			span.End()
			ctx.Abort()
			return
		}

		if !au.IsAdmin {
			err := errors.New("user is not admin. access delay")
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			zapx.WithContext(apmCtx).Error(err.Error())

			ctx.JSON(http.StatusForbidden, err)
			span.End()
			ctx.Abort()
			return
		}

		span.End()
		ctx.Next()
	}
}

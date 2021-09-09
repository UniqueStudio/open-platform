package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/UniqueStudio/open-platform/config"
	"github.com/UniqueStudio/open-platform/database"
	"github.com/UniqueStudio/open-platform/pkg"
	"github.com/UniqueStudio/open-platform/router"
	"github.com/UniqueStudio/open-platform/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/xylonx/zapx"
	zapxdecoder "github.com/xylonx/zapx/decoder"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use: "open-platform",
	PreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

var cfgFile string

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config-file", "c", "./settings.yaml", "specify config file path")
}

func Execute() error {
	return rootCmd.Execute()
}

func initConfig() {
	// init zapx
	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}

	zapx.Use(logger, zapxdecoder.OpentelemetaryDecoder)

	if err := config.SetupConfig(cfgFile); err != nil {
		zapx.Error("init config failed", zap.Error(err))
		os.Exit(1)
	}

	if err := database.SetupDatabase(); err != nil {
		zapx.Error("setup database failed", zap.Error(err))
		os.Exit(1)
	}

	if err := utils.SetupSessionStore(); err != nil {
		zapx.Error("setup session failed", zap.Error(err))
		os.Exit(1)
	}
}

func run() error {
	// setup otel tracing
	shutdown, err := utils.SetupTracing()
	defer func() {
		zapx.Info("tracing reporter is shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdown(ctx); err != nil {
			zapx.Error("tracing have been down down")
			return
		}
		zapx.Info("tracing reporter shut down successfully")
	}()

	if err != nil {
		zapx.Error("setup otel tracing failed", zap.Error(err))
		return err
	}

	if config.Config.Application.Mode != pkg.ModeDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	router.InitRouter(r)

	srv := http.Server{
		Addr:         config.Config.Application.Host + ":" + config.Config.Application.Port,
		Handler:      r,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}

	zapx.Info("start listen", zap.String("address", srv.Addr))
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			zapx.Error("listen and server failed", zap.Error(err))
			return
		}
	}()

	zapx.Info("Enter Control + C to Shutdown Server")

	// below codes are used for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	zapx.Info("server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zapx.Error("server down error: ", zap.Error(err))
		return err
	}
	zapx.Info("server have been down")
	return nil
}

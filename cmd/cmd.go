package cmd

import (
	"context"
	"log"
	"net"
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
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	if err := utils.SetupSMSClient(); err != nil {
		zapx.Error("setup tencent sms client failed", zap.Error(err))
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

	// init gin router
	r := gin.New()
	router.InitRouter(r)
	srv := http.Server{
		Addr:         config.Config.Application.Host + ":" + config.Config.Application.HttpPort,
		Handler:      r,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}

	// init grpc router
	lis, err := net.Listen("tcp", config.Config.Application.Host+":"+config.Config.Application.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := new(grpc.Server)
	if config.Config.Application.Mode == pkg.ModeDebug {
		creds, err := credentials.NewServerTLSFromFile(
			config.Config.Application.GrpcCertFile,
			config.Config.Application.GrpcKeyFile,
		)
		if err != nil {
			zapx.Fatal("new server credentials failed", zap.Error(err))
		}
		s = grpc.NewServer(
			grpc.Creds(creds),
			grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
			grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		)
	} else {
		s = grpc.NewServer(
			grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
			grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		)
	}
	router.InitGrpcHandlers(s)

	// start server
	zapx.Info("start listen tcp for http server", zap.String("port", config.Config.Application.HttpPort))
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			zapx.Error("listen and server failed", zap.Error(err))
			return
		}
	}()

	zapx.Info("start listen tcp for grpc", zap.String("port", config.Config.Application.GrpcPort))
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// below codes are used for graceful shutdown
	zapx.Info("Enter Control + C to Shutdown Server")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	zapx.Info("http server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zapx.Error("http server down error: ", zap.Error(err))
		return err
	}
	zapx.Info("http server have been down")

	zapx.Info("grpc server is shutting down...")
	s.GracefulStop()
	zapx.Info("grpc server have been down")
	return nil
}

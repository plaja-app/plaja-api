package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/plaja-app/plaja-api/pkg/logger"
	"github.com/plaja-app/plaja-api/pkg/logger/rotator"
	"github.com/plaja-app/plaja-api/user/internal/app"
	"github.com/plaja-app/plaja-api/user/internal/app/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const service = "user"

func main() {
	cfg := config.Must(config.NewFromEnv())

	l := logger.NewWithRotation(cfg.LogLevel, &rotator.Options{
		Directory: fmt.Sprintf("%s%s/", cfg.LogDirectory, service),
		MaxSize:   1,
	}).With(
		zap.String("service", service),
		zap.Int("pid", os.Getpid()),
		zap.String("address", cfg.Addr),
	)

	a := app.New(cfg, l)
	go func() {
		err := a.Run()
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			l.Error("Service exited with error", zap.Error(err))
		}
	}()

	l.Info("Service started!")

	a.Shutdown()
}

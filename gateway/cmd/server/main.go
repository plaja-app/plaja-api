package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/plaja-app/plaja-api/gateway/internal/app"

	"github.com/plaja-app/plaja-api/gateway/internal/app/config"
	"github.com/plaja-app/plaja-api/pkg/logger"
	"github.com/plaja-app/plaja-api/pkg/logger/rotator"
	"go.uber.org/zap"
)

const service = "gateway"

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
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Error("Service exited with error", zap.Error(err))
		}
	}()

	l.Info("Service started!")

	a.Shutdown()
}

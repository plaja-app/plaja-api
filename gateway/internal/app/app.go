package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/plaja-app/plaja-api/gateway/internal/app/config"
	"github.com/plaja-app/plaja-api/gateway/internal/server"
	"github.com/plaja-app/plaja-api/gateway/internal/server/mapper"
	"github.com/plaja-app/plaja-api/gateway/internal/server/marshaler"
	"github.com/plaja-app/plaja-api/pkg/logger"
	"github.com/plaja-app/plaja-api/protos/gen/go/user/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const _ = "app initialization"

const shutdownTimeout = 5 * time.Second

// New constructs new App with provided arguments.
func New(cfg *config.Config, l *logger.Logger) *App {
	return &App{
		cfg: cfg,
		l:   l,
	}
}

// App is a thin abstraction used to initialize all the dependencies.
type App struct {
	srv *http.Server
	cfg *config.Config
	l   *logger.Logger
}

func (a *App) Run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	rm := mapper.New(a.l.With(zap.String("source", "response mapper")))
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(rm.MapGRPCErr),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, marshaler.New()),
	)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := user.RegisterUserServiceHandlerFromEndpoint(ctx, mux, a.cfg.UserAddr, opts); err != nil {
		return fmt.Errorf("failed to create proxy for the gRPC service: %w", err)
	}

	s := server.New(
		middleware.Logger(mux),
		a.cfg.Addr,
		zap.NewStdLog(a.l.Get()),
	)

	a.srv = s

	a.l.Info("gRPC Gateway listening", zap.String("addr", a.cfg.Addr))
	return s.ListenAndServe()
}

// Shutdown gracefully shuts down the server. It listens to the OS signals. After
// receiving a signal, it terminates all currently active processes,  and gracefully exits.
func (a *App) Shutdown() {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	s := <-ch

	a.l.Info("Initiating shutdown...", zap.Any("signal", s))

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		a.l.Error("Failed to shutdown server", zap.Error(err))
		return
	}

	a.l.Info("Server shut down!")

	signal.Stop(ch)
}

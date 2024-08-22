package app

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/plaja-app/plaja-api/user/internal/app/config"
	"github.com/plaja-app/plaja-api/user/internal/transport/grpc/server/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/plaja-app/plaja-api/pkg/logger"
	"go.uber.org/zap"
)

const operation = "app initialization"

// New constructs new App with provided arguments.
func New(cfg *config.Config, l *logger.Logger) *App {
	return &App{
		cfg: cfg,
		l:   l,
	}
}

// App is a thin abstraction used to initialize all the dependencies.
type App struct {
	srv *grpc.Server
	cfg *config.Config
	l   *logger.Logger
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", a.cfg.Addr)
	if err != nil {
		return fmt.Errorf("%s: failed to start listener: %w", operation, err)
	}

	a.srv = grpc.NewServer()
	user.Register(a.srv, a.l)
	reflection.Register(a.srv)

	return a.srv.Serve(l)
}

// Shutdown gracefully shuts down the server. It listens to the OS signals. After
// receiving a signal, it terminates all currently active processes,  and gracefully exits.
func (a *App) Shutdown() {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	s := <-ch

	a.l.Info("Initiating shutdown...", zap.Any("signal", s))
	a.srv.Stop()
	a.l.Info("Server shut down!")

	signal.Stop(ch)
}

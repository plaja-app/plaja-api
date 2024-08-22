package user

import (
	"context"

	"github.com/plaja-app/plaja-api/pkg/logger"
	us "github.com/plaja-app/plaja-api/protos/gen/go/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const _ = "gRPC server"

// Register registers rate watcher handler to the gRPC server provided.
func Register(srv *grpc.Server, l *logger.Logger) {
	us.RegisterUserServiceServer(srv, &Server{
		l: l,
	})
}

// Server represents User Service gRPC server.
type Server struct {
	us.UnimplementedUserServiceServer
	l *logger.Logger
}

func (s *Server) Test(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	s.l.Info("Hit the user service!")
	return &emptypb.Empty{}, nil
}

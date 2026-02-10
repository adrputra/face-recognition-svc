package server

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"google.golang.org/grpc"

	"github.com/adrputra/face-recognition-svc/presence-svc/internal/config"
)

type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
}

func NewGRPCServer(cfg config.Config) (*GRPCServer, error) {
	addr := ":" + strconv.Itoa(cfg.GRPC.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("listen on %s: %w", addr, err)
	}

	grpcServer := grpc.NewServer()

	return &GRPCServer{
		server:   grpcServer,
		listener: listener,
	}, nil
}

func (s *GRPCServer) Addr() string {
	return s.listener.Addr().String()
}

func (s *GRPCServer) Start() error {
	return s.server.Serve(s.listener)
}

func (s *GRPCServer) Stop(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		s.server.Stop()
		return ctx.Err()
	}
}

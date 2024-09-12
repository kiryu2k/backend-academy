package grpcserver

import (
	"fmt"
	"homework/internal/domain"
	"homework/internal/proto"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	grpc    *grpc.Server
	handler *Handler
	port    string
}

func New(port string, file domain.FileUseCase) *server {
	grpc := grpc.NewServer(
		grpc.UnaryInterceptor(ValidateUnaryInterceptor),
		grpc.StreamInterceptor(ValidateStreamInterceptor),
	)
	s := &server{
		grpc:    grpc,
		port:    port,
		handler: NewHandler(file),
	}
	proto.RegisterFileServiceServer(grpc, s.handler)
	return s
}

func (s *server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return fmt.Errorf("listen port %s: %w", s.port, err)
	}
	if err := s.grpc.Serve(listener); err != nil {
		return fmt.Errorf("grpc serve port %s: %w", s.port, err)
	}
	return nil
}

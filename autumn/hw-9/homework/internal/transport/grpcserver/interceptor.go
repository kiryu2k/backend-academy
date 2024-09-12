package grpcserver

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type validator interface {
	Validate() error
}

func ValidateUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if v, ok := req.(validator); ok {
		if err := v.Validate(); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	return handler(ctx, req)
}

func ValidateStreamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	streamValidator := &streamValidator{
		ss: ss,
	}
	return handler(srv, streamValidator)
}

type streamValidator struct {
	ss grpc.ServerStream
}

func (s *streamValidator) RecvMsg(m any) error {
	if err := s.ss.RecvMsg(m); err != nil {
		return err
	}
	if v, ok := m.(validator); ok {
		if err := v.Validate(); err != nil {
			return status.Error(codes.InvalidArgument, err.Error())
		}
	}
	return nil
}

func (s *streamValidator) SetHeader(md metadata.MD) error {
	return s.ss.SetHeader(md)
}

func (s *streamValidator) SendHeader(md metadata.MD) error {
	return s.ss.SendHeader(md)
}

func (s *streamValidator) SetTrailer(md metadata.MD) {
	s.ss.SetTrailer(md)
}

func (s *streamValidator) Context() context.Context {
	return s.ss.Context()
}

func (s *streamValidator) SendMsg(m any) error {
	return s.ss.SendMsg(m)
}

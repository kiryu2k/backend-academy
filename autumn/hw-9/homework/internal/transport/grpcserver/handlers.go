package grpcserver

import (
	"context"
	"homework/internal/domain"
	"homework/internal/proto"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	proto.UnimplementedFileServiceServer
	file domain.FileUseCase
}

func NewHandler(file domain.FileUseCase) *Handler {
	return &Handler{
		file: file,
	}
}

func (h *Handler) Get(req *proto.GetRequest, stream proto.FileService_GetServer) error {
	file, err := h.file.Get(req.Filename)
	if err != nil {
		log.Println(err)
		return status.Error(codes.NotFound, err.Error())
	}
	err = stream.Send(&proto.GetResponse{
		File: file,
	})
	if err != nil {
		log.Println(err)
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}

func (h *Handler) All(ctx context.Context, req *proto.AllRequest) (*proto.AllResponse, error) {
	files := h.file.All(ctx)
	return &proto.AllResponse{
		Filenames: files,
	}, nil
}

func (h *Handler) GetInfo(ctx context.Context, req *proto.GetInfoRequest) (*proto.GetInfoResponse, error) {
	file, err := h.file.GetInfo(ctx, req.Filename)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &proto.GetInfoResponse{
		Filename: file.Name,
		Type:     file.Type,
		Size:     file.Size,
	}, nil
}

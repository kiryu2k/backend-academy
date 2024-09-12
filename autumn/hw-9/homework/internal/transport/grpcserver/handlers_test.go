package grpcserver_test

import (
	"context"
	"homework/internal/domain"
	"homework/internal/domain/mocks"
	"homework/internal/proto"
	"homework/internal/transport/grpcserver"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handlerSuite struct {
	suite.Suite
	usecase *mocks.FileUseCase
	handler *grpcserver.Handler
}

func (suite *handlerSuite) SetupTest() {
	suite.usecase = new(mocks.FileUseCase)
	suite.handler = grpcserver.NewHandler(suite.usecase)
}

func (suite *handlerSuite) TestGetFile() {
	testCases := []struct {
		name      string
		req       *proto.GetRequest
		resp      *proto.GetResponse
		mockRet   []byte
		err       error
		expStatus codes.Code
	}{
		{
			name: "OK",
			req: &proto.GetRequest{
				Filename: "existing_file.jpeg",
			},
			resp: &proto.GetResponse{
				File: []byte("some file data bytes"),
			},
			mockRet:   []byte("some file data bytes"),
			expStatus: codes.OK,
		},
		{
			name: "file not found",
			req: &proto.GetRequest{
				Filename: "unknown.jpeg",
			},
			err:       domain.ErrFileNotFound,
			expStatus: codes.NotFound,
		},
	}
	const methodName = "Get"
	stream := &streamMock{
		sentFromServer: make(chan *proto.GetResponse, 1),
	}
	for _, test := range testCases {
		suite.usecase.On(methodName, test.req.Filename).Return(test.mockRet, test.err)
		err := suite.handler.Get(test.req, stream)
		status := status.Code(err)
		suite.Require().Equal(test.expStatus, status, test.name)
		if test.resp != nil {
			resp, _ := stream.RecvToClient()
			suite.Require().Equal(test.resp, resp, test.name)
		}
	}
}

func (suite *handlerSuite) TestGetAllFiles() {
	const methodName = "All"
	exp := &proto.AllResponse{
		Filenames: []string{
			"file1.txt",
			"file2.jpeg",
			"totally_not_a_virus.exe",
		},
	}
	suite.usecase.On(methodName, context.Background()).Return(exp.Filenames)
	resp, err := suite.handler.All(context.Background(), new(proto.AllRequest))
	suite.Require().NoError(err)
	suite.Require().Equal(exp, resp)
}

func (suite *handlerSuite) TestGetInfoOfFile() {
	testCases := []struct {
		name      string
		req       *proto.GetInfoRequest
		resp      *proto.GetInfoResponse
		mockRet   *domain.FileInfo
		err       error
		expStatus codes.Code
	}{
		{
			name: "OK",
			req: &proto.GetInfoRequest{
				Filename: "some_file.jpeg",
			},
			resp: &proto.GetInfoResponse{
				Filename: "some_file.jpeg",
				Type:     ".jpeg",
				Size:     88,
			},
			mockRet: &domain.FileInfo{
				Name: "some_file.jpeg",
				Type: ".jpeg",
				Size: 88,
			},
			expStatus: codes.OK,
		},
		{
			name: "file not found",
			req: &proto.GetInfoRequest{
				Filename: "not_exist.jpeg",
			},
			mockRet:   nil,
			err:       domain.ErrFileNotFound,
			expStatus: codes.NotFound,
		},
	}
	const methodName = "GetInfo"
	for _, test := range testCases {
		suite.usecase.On(methodName, context.Background(), test.req.Filename).Return(test.mockRet, test.err)
		resp, err := suite.handler.GetInfo(context.Background(), test.req)
		status := status.Code(err)
		suite.Require().Equal(test.expStatus, status, test.name)
		suite.Require().Equal(test.resp, resp, test.name)
	}
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(handlerSuite))
}

type streamMock struct {
	grpc.ServerStream
	sentFromServer chan *proto.GetResponse
}

func (s *streamMock) Context() context.Context {
	return context.Background()
}

func (s *streamMock) Send(resp *proto.GetResponse) error {
	s.sentFromServer <- resp
	return nil
}

func (s *streamMock) RecvToClient() (*proto.GetResponse, error) {
	return <-s.sentFromServer, nil
}

func (s *streamMock) Recv() (*proto.GetRequest, error) {
	return nil, nil
}
func (s *streamMock) SendFromClient(req *proto.GetRequest) error {
	return nil
}

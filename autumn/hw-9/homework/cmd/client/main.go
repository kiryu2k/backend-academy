package main

import (
	"context"
	"flag"
	"homework/internal/proto"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	GetMethod     = "get"
	AllMethod     = "all"
	GetInfoMethod = "getinfo"
)

func main() {
	var (
		method   = flag.String("m", GetMethod, "method to call")
		port     = flag.Int("p", 50051, "port of grpc server")
		filename = flag.String("f", "", "file name to find")
	)
	flag.Parse()
	conn, err := grpc.Dial(getAddress(*port), grpc.WithTransportCredentials(
		insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	cli := proto.NewFileServiceClient(conn)
	switch strings.ToLower(*method) {
	case GetMethod:
		getFilesRequest(cli, *filename)
	case AllMethod:
		getAllFilesRequest(cli)
	case GetInfoMethod:
		getFileInfoRequest(cli, *filename)
	default:
		log.Fatalf("unsupported method: %s", *method)
	}
}

func getFilesRequest(cli proto.FileServiceClient, filename string) {
	stream, err := cli.Get(context.Background(), &proto.GetRequest{
		Filename: filename,
	})
	if err != nil {
		log.Fatalf("error on stream messages: %v", err)
	}
	for {
		file, err := stream.Recv()
		if err == nil {
			log.Printf("received file: %s\n\n", file)
			continue
		}
		if err == io.EOF {
			return // end of stream
		}
		log.Fatalf("error while receiving file: %v", err)
	}
}

func getAllFilesRequest(cli proto.FileServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := cli.All(ctx, &proto.AllRequest{})
	if err != nil {
		log.Fatalf("failed to get all files: %v", err)
	}
	log.Printf("received files: %v\n", resp.Filenames)
}

func getFileInfoRequest(cli proto.FileServiceClient, filename string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := cli.GetInfo(ctx, &proto.GetInfoRequest{
		Filename: filename,
	})
	if err != nil {
		log.Fatalf("failed to get file info: %v", err)
	}
	log.Printf("received file info: %v\n", resp)
}

func getAddress(port int) string {
	return ":" + strconv.Itoa(port)
}

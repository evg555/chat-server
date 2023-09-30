package main

import (
	"context"
	"log"
	"net"
	"strings"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	service "github.com/evg555/chat-server/pkg/chat_v1"
)

const address = "localhost:8000"

type server struct {
	service.UnimplementedChatServer
}

func main() {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	service.RegisterChatServer(s, &server{})

	log.Printf("server is starting at %s", address)

	if err = s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) Create(_ context.Context, req *service.CreateRequest) (*service.CreateResponse, error) {
	log.Println("creating chat...")
	log.Printf("users: %s", strings.Join(req.GetUsernames(), ", "))

	return &service.CreateResponse{Id: gofakeit.Int64()}, nil
}

func (s *server) Delete(_ context.Context, req *service.DeleteRequest) (*emptypb.Empty, error) {
	log.Println("deleting chat...")
	log.Printf("chat id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(_ context.Context, req *service.SendMessageRequest) (*emptypb.Empty, error) {
	log.Println("send message to chat...")
	log.Printf("message: %+v", req)

	return &emptypb.Empty{}, nil
}

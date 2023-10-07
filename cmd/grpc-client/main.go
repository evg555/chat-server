package main

import (
	"context"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	service "github.com/evg555/chat-server/pkg/chat_v1"
)

const (
	address = "localhost:8000"
	chatId  = 42
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c := service.NewChatClient(conn)

	resp1, err := c.Create(ctx, &service.CreateRequest{Usernames: []string{
		gofakeit.Username(),
		gofakeit.Username(),
		gofakeit.Username(),
	}})
	if err != nil {
		log.Fatalf("failed to create chat: %v", err)
	}

	log.Printf(color.RedString("Create chat\n"), color.GreenString("chat id: %d", resp1.GetId()))

	resp2, err := c.Delete(ctx, &service.DeleteRequest{Id: chatId})
	if err != nil {
		log.Fatalf("failed to delete chat: %v", err)
	}

	log.Printf(color.RedString("Delete chat\n"), color.GreenString("resp: %+v", resp2))

	resp3, err := c.SendMessage(ctx, &service.SendMessageRequest{
		From:      gofakeit.Username(),
		Text:      gofakeit.HackerPhrase(),
		Timestamp: timestamppb.New(gofakeit.Date()),
	})
	if err != nil {
		log.Fatalf("failed to send message to chat: %v", err)
	}

	log.Printf(color.RedString("Send message to chat\n"), color.GreenString("resp: %+v", resp3))
}

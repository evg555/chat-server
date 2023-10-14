package main

import (
	"context"
	"log"
	"net"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	service "github.com/evg555/chat-server/pkg/chat_v1"
)

const (
	address = "localhost:8000"
	dbDSN   = "host=localhost port=5432 dbname=chat user=chat password=chat sslmode=disable"
	table   = "chat"
)

type server struct {
	service.UnimplementedChatServer
	db *pgx.Conn
}

func main() {
	ctx := context.Background()

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	conn, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	s := grpc.NewServer()
	reflection.Register(s)
	service.RegisterChatServer(s, &server{db: conn})

	log.Printf("server is starting at %s", address)

	if err = s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) Create(ctx context.Context, req *service.CreateRequest) (*service.CreateResponse, error) {
	log.Println("creating chat...")
	log.Printf("users: %s", strings.Join(req.GetUsernames(), ", "))

	var id int64

	builderInsert := sq.Insert(table).
		PlaceholderFormat(sq.Dollar).
		Columns("user_from", "user_to").
		Values(req.GetUsernames()[0], req.GetUsernames()[1]).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("failed to insert to database: %v", err)
		return nil, err
	}

	res := s.db.QueryRow(ctx, query, args...)

	err = res.Scan(&id)
	if err != nil {
		log.Printf("failed to insert to database: %v", err)
		return nil, err
	}

	return &service.CreateResponse{Id: id}, nil
}

func (s *server) Delete(ctx context.Context, req *service.DeleteRequest) (*emptypb.Empty, error) {
	log.Println("deleting chat...")
	log.Printf("chat id: %d", req.GetId())

	builderDelete := sq.Delete(table).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Printf("failed to delete from database: %v", err)
		return nil, err
	}

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to delete from database: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *service.SendMessageRequest) (*emptypb.Empty, error) {
	log.Println("send message to chat...")
	log.Printf("message: %+v", req)

	builderDelete := sq.Update(table).
		PlaceholderFormat(sq.Dollar).
		Set("user_from", req.GetFrom()).
		Set("text", req.GetText()).
		Set("timestamp", req.GetTimestamp().AsTime()).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Printf("failed to update to database: %v", err)
		return nil, err
	}

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update to database: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

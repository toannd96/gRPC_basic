package main

import (
	"context"
	"fmt"
	"grpc-learn/platform/platformpb"
	"grpc-learn/platform/server/providers"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	defaultPort = 3000
)

type server struct {
	bigCommerce *providers.BigCommerce
	magento     *providers.Magento
}

func (s *server) GetProduct(ctx context.Context, req *platformpb.ProductRequest) (*platformpb.ProductResponse, error) {
	log.Println("server run...", req.Platform)

	response, err := s.bigCommerce.Query(req.Platform)
	if err != nil {
		return nil, err
	}

	return &platformpb.ProductResponse{
		Name:       response.Name,
		Price:      response.Price,
		Sku:        response.Sku,
		Type:       response.Type,
		Categories: response.Categories,
	}, nil
}

func main() {
	conn := listen()
	grpcServer := grpc.NewServer()
	platformpb.RegisterProductServiceServer(grpcServer, &server{})
	grpcServer.Serve(conn)
}

func listen() net.Listener {
	listenAddr := fmt.Sprintf(":%d", defaultPort)

	conn, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("err while create listen %v", err)
	}
	log.Println("[ProductServer] Listening on", defaultPort)
	return conn
}

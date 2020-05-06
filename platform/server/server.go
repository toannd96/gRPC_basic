package main

import (
	"context"
	"fmt"
	"grpc-learn/platform/platformpb"
	"grpc-learn/platform/server/providers"
	"log"
	"net"

	"github.com/joho/godotenv"

	"google.golang.org/grpc"
)

const (
	port = 3000
)

var (
	accessTokenBigCommerce string
	accessTokenMagento     string
	clientID               string
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("not receive environment variables")
	}
}

type Server struct {
	bigCommerce *providers.BigCommerce
	magento     *providers.Magento
}

func (s *Server) GetProduct(ctx context.Context, req *platformpb.ProductRequest) (*platformpb.ProductResponse, error) {
	log.Println("[server run] get product from platform", req.Platform)

	providerName := req.GetPlatform()

	var listResponse []providers.PlatformInfo

	if providerName == s.bigCommerce.Name() {
		bigcommerceAPI := "https://api.bigcommerce.com/stores/jp6bmqxeb9/v3/catalog/products"
		response, err := s.bigCommerce.Get(bigcommerceAPI)
		if err != nil {
			return nil, err
		}
		listResponse = append(listResponse, response)
	}

	if providerName == s.magento.Name() {
		magentoAPI := "https://magento23demo.connectpos.com/rest/V1/products?searchCriteria[pageSize]=1"
		response, err := s.magento.Get(magentoAPI)
		if err != nil {
			return nil, err
		}
		listResponse = append(listResponse, response)
	}
	return &platformpb.ProductResponse{
		Name:       listResponse[0].Name,
		Price:      listResponse[0].Price,
		Sku:        listResponse[0].Sku,
		Type:       listResponse[0].Type,
		Categories: listResponse[0].Categories,
	}, nil
}

func main() {
	conn := listen()
	// khởi tạo một đối tượng gRPC service
	grpcServer := grpc.NewServer()

	// đăng ký service với grpcServer
	platformpb.RegisterProductServiceServer(grpcServer, &Server{})
	grpcServer.Serve(conn)
}

func listen() net.Listener {
	listenAddr := fmt.Sprintf(":%d", port)

	// cung cấp gRPC service trên port
	conn, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("err while create listen %v", err)
	}
	log.Println("[PlatformServer] Listening on", port)
	return conn
}

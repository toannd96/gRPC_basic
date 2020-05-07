package main

import (
	"context"
	"log"

	"grpc-learn/platform/platformpb"

	"google.golang.org/grpc"
)

const (
	bigcommerceName = "bigcommerce"
	magentoName     = "magento"
)

func main() {
	// thiết lập kết nối với gRPC service
	conn, err := grpc.Dial("localhost:3000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// xây dựng đối tượng `ProductServiceClient` dựa trên kết nối đã thiết lập
	client := platformpb.NewProductServiceClient(conn)
	callGetProduct(client)
}

func callGetProduct(p platformpb.ProductServiceClient) {
	log.Println("[client run] calling get product api...")
	req := platformpb.ProductRequest{
		// Platform: bigcommerceName,
		// Name:     "Configurable_mix full options_1001",
		Platform: magentoName,
		Name:     "Bag and strap",
	}

	resp, err := p.GetProduct(context.Background(), &req)
	if err != nil {
		log.Fatalf("request failed: %v", err)
	}

	if resp.GetName() == "" {
		log.Fatal("product not found")
	}
	log.Println(resp)
}

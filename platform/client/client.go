package main

import (
	"context"
	"log"

	"grpc-learn/platform/platformpb"

	"google.golang.org/grpc"
)

const (
	name1 = "bigcommerce"
	name2 = "magento23demo"
)

func main() {
	conn, err := grpc.Dial("localhost:3000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := platformpb.NewProductServiceClient(conn)
	callGetProduct(client)
}

func callGetProduct(p platformpb.ProductServiceClient) {
	log.Println("calling get product api...")

	// api := "https://magento23demo.connectpos.com/rest/V1/products?searchCriteria[pageSize]=1"

	req := platformpb.ProductRequest{
		Platform: name1,
	}
	resp, err := p.GetProduct(context.Background(), &req)
	if err != nil {
		log.Fatalf("request failed: %v", err)
	}
	log.Println(resp)
}

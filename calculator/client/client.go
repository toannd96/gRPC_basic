package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"

	"grpc-learn/calculator/calculatorpb"
)

func main() {
	cc, err := grpc.Dial("localhost:50069", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("err while dial %v", err)
	}

	defer cc.Close()

	client := calculatorpb.NewCalculatorServiceClient(cc)

	// log.Printf("service client %f", client)
	// callSum(client)
	// callPND(client)
	// callAverage(client)
	callFindMax(client)
}

func callSum(c calculatorpb.CalculatorServiceClient) {
	log.Println("calling sum api")
	resp, err := c.Sum(context.Background(), &calculatorpb.SumRequest{
		Number1: 7,
		Number2: 6,
	})

	if err != nil {
		log.Fatalf("call sum api err %v", err)
	}

	log.Printf("sum api response %v\n", resp.GetResult())
}

func callPND(c calculatorpb.CalculatorServiceClient) {
	log.Println("calling pnd api")
	stream, err := c.PrimeNumberDecomposition(context.Background(), &calculatorpb.PNDRequest{
		Number: 300,
	})

	if err != nil {
		log.Fatalf("callPND err %v", err)
	}

	for {
		resp, recvErr := stream.Recv()

		// stream finish
		if recvErr == io.EOF {
			log.Println("server finish streaming")
			return
		}

		log.Printf("prime number %v", resp.GetResult())
	}
}

func callAverage(c calculatorpb.CalculatorServiceClient) {
	log.Println("calling average api")
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("call average err %v", err)
	}

	listReq := []calculatorpb.AverageRequest{
		calculatorpb.AverageRequest{
			Number: 5,
		},
		calculatorpb.AverageRequest{
			Number: 10,
		},
		calculatorpb.AverageRequest{
			Number: 15.3,
		},
		calculatorpb.AverageRequest{
			Number: 2.4,
		},
		calculatorpb.AverageRequest{
			Number: 6.7,
		},
	}

	for _, req := range listReq {
		err := stream.Send(&req)
		if err != nil {
			log.Fatalf("send average request err %v", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("receive average response err %v", err)
	}

	log.Printf("average response %v", resp)
}

func callFindMax(c calculatorpb.CalculatorServiceClient) {
	log.Println("calling find max api")
	stream, err := c.FindMax(context.Background())
	if err != nil {
		log.Fatalf("call find max err %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		// send many requests
		listReq := []calculatorpb.FindMaxRequest{
			calculatorpb.FindMaxRequest{
				Number: 5,
			},
			calculatorpb.FindMaxRequest{
				Number: 10,
			},
			calculatorpb.FindMaxRequest{
				Number: 3,
			},
			calculatorpb.FindMaxRequest{
				Number: 2,
			},
			calculatorpb.FindMaxRequest{
				Number: 7,
			},
		}
		for _, req := range listReq {
			err := stream.Send(&req)
			if err != nil {
				log.Fatalf("send finmax request err %v", err)
				break
			}
			time.Sleep(1000 * time.Millisecond)
		}
		// client finish
		stream.CloseSend()
	}()

	go func() {
		// receive many requests
		for {
			resp, err := stream.Recv()

			// server finish
			if err == io.EOF {
				log.Println("ending find max api")
				break
			}
			if err != nil {
				log.Fatalf("receive find max err %v", err)
				break
			}
			log.Printf("max: %v\n", resp.GetMax())
		}
		close(waitc)
	}()

	<-waitc
}

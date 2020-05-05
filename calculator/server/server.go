package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"grpc-learn/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("sum called...")
	resp := &calculatorpb.SumResponse{
		Result: req.GetNumber1() + req.GetNumber2(),
	}

	return resp, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PNDRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	log.Println("PrimeNumberDecomposition called...")
	k := int32(2)
	number := req.GetNumber()
	for number > 1 {
		if number%k == 0 {
			number = number / k

			// send to client
			stream.Send(&calculatorpb.PNDResponse{
				Result: k,
			})
		} else {
			k++
			log.Printf("k increase to %v", k)
		}
	}
	return nil
}

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	log.Println("Average called...")
	var total float32
	var count int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			resp := &calculatorpb.AverageResponse{
				Result: total / float32(count),
			}

			return stream.SendAndClose(resp)
		}
		if err != nil {
			log.Fatalf("err while receive average %v", err)
			return err
		}
		log.Printf("receive req %v", req)
		total += req.GetNumber()
		count++
	}
}

func (*server) FindMax(stream calculatorpb.CalculatorService_FindMaxServer) error {
	log.Println("Find max called...")
	max := int32(0)
	for {
		req, err := stream.Recv()

		// client finish
		if err == io.EOF {
			log.Println("EOF...")
			return nil
		}
		if err != nil {
			log.Fatalf("err while receive find max %v", err)
			return err
		}

		number := req.GetNumber()
		log.Printf("receive number %v\n", number)
		if number > max {
			max = number
		}

		// send to client
		err = stream.Send(&calculatorpb.FindMaxResponse{
			Max: max,
		})
		if err != nil {
			log.Fatalf("send max err %v", err)
			return err
		}
	}
}
func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50069")
	if err != nil {
		log.Fatalf("err while create listen %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	fmt.Println("calculator is running...")
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("err while server %v", err)
	}
}

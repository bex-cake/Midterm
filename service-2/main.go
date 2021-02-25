package main

import (
	"awesomeProject4/service-2/proto/calculatorpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	fmt.Println("Hello I'm a client")
	HOST := os.Getenv("HOST")
	PORT := os.Getenv("PORT")
	URL := HOST + ":" + PORT

	conn, err := grpc.Dial(URL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)

	doPrimeDecomposition(c)
	doComputeAverage(c)
}

func doPrimeDecomposition(c calculatorpb.CalculatorServiceClient) {
	ctx := context.Background()
	request := &calculatorpb.PrimeDecompositionRequest{
		Number: &calculatorpb.Number{Number: 210},
	}

	stream, err := c.PrimeDecomposition(ctx, request)

	if err != nil {
		log.Fatalf("Error while calling PrimeDecomposition RPC %v", err)
	}
	defer stream.CloseSend()

	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Error while reciving from PrimeDecomposition RPC %v", err)
		}
		log.Printf("Response from PrimeDecomposition RPC:%v \n", res.GetResult())
	}
}

func doComputeAverage(c calculatorpb.CalculatorServiceClient) {

	requests := []*calculatorpb.ComputeAverageRequest{
		{
			Number: &calculatorpb.Number{Number: 5},
		},
		{
			Number: &calculatorpb.Number{Number: 7},
		},
		{
			Number: &calculatorpb.Number{Number: 13},
		},
	}

	ctx := context.Background()
	stream, err := c.ComputeAverage(ctx)
	if err != nil {
		log.Fatalf("Error while calling ComputeAverage: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from ComputeAverage: %v", err)
	}
	fmt.Printf("ComputeAverage Response: %v\n", res)
}

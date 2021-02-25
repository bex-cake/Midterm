package main

import (
	"awesomeProject4/service-1/proto/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type Server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (s *Server) PrimeDecomposition(req *calculatorpb.PrimeDecompositionRequest, stream calculatorpb.CalculatorService_PrimeDecompositionServer) error {
	number := req.GetNumber().GetNumber()
	k := int32(2)
	for number > 1 {
		if number%k == 0 {
			res := &calculatorpb.PrimeDecompositionResponse{Result: k}
			number = number / k
			if err := stream.Send(res); err != nil {
				log.Fatalf("Error while sending PrimeDecomposition requests: %v", err.Error())
			}
			time.Sleep(time.Second)
		} else {
			k = k + 1
		}
	}
	return nil
}

func (s *Server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	var total int32
	var count float64
	total = 0
	count = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{Result: float64(total) / count})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		count++
		total += req.Number.GetNumber()
	}
}

func main() {
	PORT := os.Getenv("PORT")
	l, err := net.Listen("tcp", "0.0.0.0:"+PORT)
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &Server{})
	log.Println("Server is running on port:" + PORT)
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}
}

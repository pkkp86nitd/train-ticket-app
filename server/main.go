package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/pkkp86nitd/train_ticket_app/proto"
	customserver "github.com/pkkp86nitd/train_ticket_app/server/customServer"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	server := customserver.GetCustomServerInstance()
	pb.RegisterTrainTicketServiceServer(s, server)
	fmt.Println("Server is listening on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

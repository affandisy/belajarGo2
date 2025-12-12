package main

import (
	"log"
	"net"

	pb "belajarGo2/app/grpc-server/controller/inventory"

	"google.golang.org/grpc"
)

func main() {
	// Listen on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the service implementation
	pb.RegisterInventoryServiceServer(grpcServer, pb.NewInventoryService())

	log.Println("gRPC server running on port 50051...")

	// Start serving
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

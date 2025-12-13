package main

import (
	"log"
	"log/slog"
	"net"
	"os"
	"strings"

	pb "belajarGo2/app/grpc-server/controller/inventory"
	"belajarGo2/app/grpc-server/middleware"

	cfg "github.com/pobyzaarif/go-config"

	"google.golang.org/grpc"
)

var loggerOption = slog.HandlerOptions{AddSource: true}
var logger = slog.New(slog.NewJSONHandler(os.Stdout, &loggerOption))

type Config struct {
	AppPort      string `env:"APP_PORT_GRPC_SERVER"`
	AppBasicAuth string `env:"APP_BASIC_AUTH"`
}

func main() {
	config := Config{}
	cfg.LoadConfig(&config)
	logger.Info("Config loaded")

	// Extract basic auth credentials from config
	basicAuthMap := make(map[string]string)
	basicAuthConfig := strings.Split(config.AppBasicAuth, ",")
	for _, v := range basicAuthConfig {
		basicAuthPair := strings.Split(v, ":")
		basicAuthMap[basicAuthPair[0]] = basicAuthPair[1]
	}

	// Listen grpc with port from config
	lis, err := net.Listen("tcp", ":"+config.AppPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	// grpcServer := grpc.NewServer()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.BasicAuthUnaryInterceptor(basicAuthMap)),
		grpc.StreamInterceptor(middleware.BasicAuthStreamInterceptor(basicAuthMap)),
	)

	// Register the service implementation
	pb.RegisterInventoryServiceServer(grpcServer, pb.NewInventoryService())

	// log.Println("gRPC server running on port 50051...")
	logger.Info("gRPC service running on port " + config.AppPort)

	// Start serving
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

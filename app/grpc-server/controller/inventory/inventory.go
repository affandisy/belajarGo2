package inventory

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type inventoryServiceServer struct {
	UnimplementedInventoryServiceServer
}

func NewInventoryService() InventoryServiceServer {
	return &inventoryServiceServer{}
}

func (s *inventoryServiceServer) Create(ctx context.Context, req *InventoryRequest) (*InventoryResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "request is nil")
	}
	if req.GetCode() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "code is required")
	}
	if req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}
	if req.GetStatus() == InventoryStatus_INVENTORY_STATUS_UNSPECIFIED {
		return nil, status.Errorf(codes.InvalidArgument, "status is required")
	}

	resp := &InventoryResponse{Inventory: req}
	return resp, nil
}

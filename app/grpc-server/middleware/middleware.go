package middleware

import (
	"context"
	"encoding/base64"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func BasicAuthUnaryInterceptor(allowed map[string]string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if len(allowed) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
		}
		if ok := checkBasicAuthAgainstMap(ctx, allowed); !ok {
			return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
		}
		return handler(ctx, req)
	}
}

func BasicAuthStreamInterceptor(allowed map[string]string) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if len(allowed) == 0 {
			return status.Errorf(codes.Unauthenticated, "unauthenticated")
		}
		if ok := checkBasicAuthAgainstMap(ss.Context(), allowed); !ok {
			return status.Errorf(codes.Unauthenticated, "unauthenticated")
		}
		return handler(srv, ss)
	}
}

func checkBasicAuthAgainstMap(ctx context.Context, allowed map[string]string) bool {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false
	}
	vals := md.Get("authorization")
	if len(vals) == 0 {
		return false
	}
	auth := vals[0]
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 {
		return false
	}
	scheme := strings.ToLower(strings.TrimSpace(parts[0]))
	if scheme != "basic" {
		return false
	}
	payload := strings.TrimSpace(parts[1])
	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return false
	}
	pair := string(decoded)
	up := strings.SplitN(pair, ":", 2)
	if len(up) != 2 {
		return false
	}
	user := up[0]
	pass := up[1]
	if allowedPass, exists := allowed[user]; exists && allowedPass == pass {
		return true
	}
	return false
}

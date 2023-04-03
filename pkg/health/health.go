package health

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var ErrServiceError = errors.New("service NOT SERVING")

// Check gRPC协议的健康检查
func Check(ctx context.Context, addr string, creds credentials.TransportCredentials) error {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	healthClient := grpc_health_v1.NewHealthClient(conn)
	healthCheckReq := &grpc_health_v1.HealthCheckRequest{}

	resp, err := healthClient.Check(ctx, healthCheckReq)
	if err != nil {
		return err
	}

	if resp.Status == grpc_health_v1.HealthCheckResponse_SERVING {
		return nil
	} else {
		return ErrServiceError
	}
}

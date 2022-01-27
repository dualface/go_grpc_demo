package impl

import (
	"context"
	"fmt"
	"go_grpc_demo/server"
	rpc "go_grpc_demo/server/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DemoServer struct {
	rpc.UnimplementedDemoServer
	Hub server.Hub
}

func (s DemoServer) Say(ctx context.Context, req *rpc.SayRequest) (*rpc.SayReply, error) {
	client, ok := s.Hub.GetAuthedClientByContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "not found authed client")
	}

	return &rpc.SayReply{Message: fmt.Sprintf("Hello, CLIENT[%d] %s", client.Id(), req.Content)}, nil
}

package impl

import (
	"context"
	"go_grpc_demo/server"
	rpc "go_grpc_demo/server/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type AuthServer struct {
	rpc.UnimplementedAuthServer
	Hub server.Hub
}

func (s AuthServer) Auth(ctx context.Context, req *rpc.AuthRequest) (*rpc.AuthReply, error) {
	client, ok := s.Hub.GetClientByContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "not found client")
	}

	if len(req.Token) > 0 {
		client.SetAuthed()
		log.Printf("[Auth]: client '%d' auth ok\n", client.Id())
		return &rpc.AuthReply{Ok: 1}, nil
	} else {
		log.Printf("[Auth]: client '%d' auth failed\n", client.Id())
		return nil, status.Errorf(codes.PermissionDenied, "invalid token")
	}
}

package impl

import (
	"context"
	"go_grpc_demo/server"
	"google.golang.org/grpc/stats"
	"log"
)

type (
	gRPCStatsHandler struct {
		hub *hub
	}
)

func (h gRPCStatsHandler) TagRPC(ctx context.Context, _ *stats.RPCTagInfo) context.Context {
	return ctx
}

func (h gRPCStatsHandler) HandleRPC(context.Context, stats.RPCStats) {
}

func (h gRPCStatsHandler) TagConn(ctx context.Context, s *stats.ConnTagInfo) context.Context {
	addr := s.RemoteAddr.String()
	client, err := h.hub.newClient(ctx, addr)
	if err != nil {
		log.Panicf("[HUB]: create client failed, %s", err)
	}
	// 在 HUB 中注册客户端
	h.hub.register <- client
	return client.ctx
}

func (h gRPCStatsHandler) HandleConn(ctx context.Context, s stats.ConnStats) {
	if _, ok := s.(*stats.ConnEnd); ok {
		id, ok := server.GetClientIdInContext(ctx)
		if !ok {
			log.Panicf("[HUB]: not found ClientId in context")
		}
		client, ok := h.hub.getClient(id)
		if !ok {
			log.Panicf("[HUB]: not found client '%d' in Hub", id)
		}
		// 从 HUB 注销客户端
		h.hub.unregister <- client
	}
}

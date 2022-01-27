package impl

import (
	"context"
	"fmt"
	"go_grpc_demo/config"
	"go_grpc_demo/server"
	rpc "go_grpc_demo/server/rpc"
	"go_grpc_demo/server/rpc/impl"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type (
	hub struct {
		cfg        *config.GlobalConfig
		clients    map[int64]*client
		register   chan *client
		unregister chan *client
		// 用于跟踪客户端连接
		conns        map[string]net.Conn
		nextClientId int64
	}
)

func NewHub(cfg *config.GlobalConfig) *hub {
	h := &hub{
		cfg:        cfg,
		clients:    make(map[int64]*client),
		register:   make(chan *client),
		unregister: make(chan *client),
		conns:      make(map[string]net.Conn),
	}

	return h
}

func (h *hub) Run() {
	// 监听网络端口
	lis, err := net.Listen("tcp", h.cfg.Server.Listen)
	if err != nil {
		log.Fatalf("[HUB]: failed to listen '%s', %s", h.cfg.Server.Listen, err)
	}

	// 创建 gRPC 服务器
	var opts []grpc.ServerOption
	// gRPCTransportCredentials 配合 gRPCStatsHandler 完成对客户端连接的跟踪
	opts = append(opts, grpc.Creds(&gRPCTransportCredentials{hub: h}))
	opts = append(opts, grpc.StatsHandler(&gRPCStatsHandler{hub: h}))
	s := grpc.NewServer(opts...)

	// 注册 gRPC 服务
	rpc.RegisterAuthServer(s, impl.AuthServer{Hub: h})
	rpc.RegisterDemoServer(s, impl.DemoServer{Hub: h})
	log.Printf("[HUB]: server listening at %v", lis.Addr())

	// 启动管理器，处理客户端注册
	h.startClientManager()

	// 启动 gRPC 服务器
	if err := s.Serve(lis); err != nil {
		log.Fatalf("[HUB]: failed to serve: %v", err)
	}
}

func (h *hub) GlobalConfig() *config.GlobalConfig {
	return h.cfg
}

func (h *hub) NumberOfClients() int {
	return len(h.clients)
}

func (h *hub) GetClientByContext(ctx context.Context) (server.Client, bool) {
	return h.getClientByContext(ctx)
}

func (h *hub) GetAuthedClientByContext(ctx context.Context) (server.Client, bool) {
	client, ok := h.GetClientByContext(ctx)
	if !ok || !client.Auth() {
		return nil, false
	}
	return client, true
}

//// private

func (h *hub) newClient(ctx context.Context, addr string) (*client, error) {
	conn, ok := h.conns[addr]
	delete(h.conns, addr)
	if !ok {
		return nil, fmt.Errorf("not found conn with addr '%s", addr)
	}
	h.nextClientId++
	id := h.nextClientId
	client := &client{id: id, hub: h, ctx: server.NewContextWithClientId(ctx, id), conn: conn}
	log.Printf("[HUB]: create client '%d', remote addr: %s", id, addr)
	client.waitForAuth(h.cfg.Server.AuthTimeout * time.Millisecond)
	return client, nil
}

func (h *hub) getClient(cid int64) (*client, bool) {
	c, ok := h.clients[cid]
	return c, ok
}

func (h *hub) getClientByContext(ctx context.Context) (*client, bool) {
	id, ok := server.GetClientIdInContext(ctx)
	if !ok {
		return nil, false
	}
	return h.getClient(id)
}

func (h *hub) startClientManager() {
	go func() {
		for {
			select {
			case client := <-h.register:
				h.clients[client.id] = client
				log.Printf("[HUB]: CLIENT[%d] registered", client.id)

			case client := <-h.unregister:
				_, ok := h.clients[client.id]
				if !ok {
					// 客户端未登记
					log.Printf("[HUB]: try to remove unregistered CLIENT[%d]", client.id)
				} else {
					delete(h.clients, client.id)
					log.Printf("[HUB]: CLIENT[%d] unregistered", client.id)
				}
			}
		}
	}()

	log.Printf("[HUB]: waiting for clients")
}

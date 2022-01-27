package impl

import (
	"context"
	"google.golang.org/grpc/credentials"
	"net"
)

type (
	gRPCTransportCredentials struct {
		hub *hub
	}

	authType struct {
		credentials.CommonAuthInfo
	}
)

func (a *authType) AuthType() string {
	return "auth"
}

func (tc *gRPCTransportCredentials) ClientHandshake(_ context.Context, _ string, c net.Conn) (net.Conn, credentials.AuthInfo, error) {
	return c, &authType{}, nil
}

func (tc *gRPCTransportCredentials) ServerHandshake(c net.Conn) (net.Conn, credentials.AuthInfo, error) {
	// 使用 RemoteAddr() 关联 Conn 和后面即将创建的 Client
	addr := c.RemoteAddr().String()
	tc.hub.conns[addr] = c
	return c, &authType{}, nil
}

func (tc *gRPCTransportCredentials) Info() credentials.ProtocolInfo {
	return credentials.ProtocolInfo{}
}

func (tc *gRPCTransportCredentials) Clone() credentials.TransportCredentials {
	return &gRPCTransportCredentials{}
}

func (tc *gRPCTransportCredentials) OverrideServerName(string) error {
	return nil
}

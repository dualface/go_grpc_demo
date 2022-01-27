package server

import (
	"context"
)

type (
	Client interface {
		Id() int64

		// Hub 返回客户端所属的 Hub
		Hub() Hub

		// Context 返回连接绑定的 Context
		Context() context.Context

		// SetAuthed 设置客户端已经验证通过
		SetAuthed()

		// Auth 返回客户端验证状态
		Auth() bool
	}

	ClientIdKeyStruct struct {
	}
)

var (
	ClientIdKey = &ClientIdKeyStruct{}
)

func NewContextWithClientId(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, ClientIdKey, id)
}

func GetClientIdInContext(ctx context.Context) (int64, bool) {
	v := ctx.Value(ClientIdKey)
	id, ok := v.(int64)
	if !ok {
		return 0, false
	}
	return id, true
}

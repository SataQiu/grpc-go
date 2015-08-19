package monitoring

import (
	"google.golang.org/grpc/codes"
)

type RpcType string

const (
	Unary     RpcType = "unary"
	Streaming RpcType = "streaming"
)

// RpcMonitor is a per-RPC datastructure
type RpcMonitor interface {
	// ReceivedMessage is called on every stream message received by the monitor.
	ReceivedMessage()

	// SentMessage is called on every stream message sent by the monitor.
	SentMessage()

	// Handled is called whenever the RPC handling completes (with OK or AppError).
	Handled(code codes.Code)

	// Errored is called whenever the RPC failed due to RPC-layer errors.
	Erred(err error)
}

// ServerMonitor allocates new per-RPC monitors on the server side.
type ServerMonitor interface {
	// NewMonitor allocates a new per-RPC monitor, also signifying the start of an RPC call.
	NewServerMonitor(rpcType RpcType, fullMethod string) RpcMonitor
}

// ServerMonitor allocates new per-RPC monitors on the server side.
type ClientMonitor interface {
	// NewMonitor allocates a new per-RPC monitor, also signifying the start of an RPC call.
	NewClientMonitor(rpcType RpcType, fullMethod string) RpcMonitor
}

// NoOpMonitor is both a Client- and ServerMonitor that does nothing.
type NoOpMonitor struct{}

func (m *NoOpMonitor) NewServerMonitor(RpcType, string) RpcMonitor {
	return m
}

func (m *NoOpMonitor) NewClientMonitor(RpcType, string) RpcMonitor {
	return m
}

func (*NoOpMonitor) ReceivedMessage() {}

func (*NoOpMonitor) SentMessage() {}

func (*NoOpMonitor) Handled(code codes.Code) {}

func (*NoOpMonitor) Erred(err error) {}

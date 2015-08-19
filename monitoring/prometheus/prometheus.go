package prometheus

import (
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/monitoring"
	"google.golang.org/grpc/transport"

	prom "github.com/prometheus/client_golang/prometheus"
)

var (
	serverStartedCounter = prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: "grpc",
			Subsystem: "server",
			Name:      "rpc_started_total",
			Help:      "Total number of RPCs started by the server.",
		}, []string{"type", "service", "method"})

	serverStreamMsgReceived = prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: "grpc",
			Subsystem: "server",
			Name:      "rpc_msg_received_total",
			Help:      "Total number of RPC stream messages received on the server.",
		}, []string{"type", "service", "method"})

	serverStreamMsgSent = prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: "grpc",
			Subsystem: "server",
			Name:      "rpc_msg_sent_total",
			Help:      "Total number of RPC stream messages sent by the server.",
		}, []string{"type", "service", "method"})

	serverHandledHistogram = prom.NewHistogramVec(
		prom.HistogramOpts{
			Namespace: "grpc",
			Subsystem: "server",
			Name:      "rpc_handled",
			Help:      "Histogram of response latency of RPC that had been application-level handled by the server.",
			Buckets:   prom.DefBuckets,
		}, []string{"type", "service", "method", "code"})

	serverErred = prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: "grpc",
			Subsystem: "server",
			Name:      "rpc_erred_total",
			Help:      "Total number of RPC that had failed on the RPC layer on the server.",
		}, []string{"type", "service", "method", "error"})
)

func init() {
	prom.MustRegister(serverStartedCounter)
	prom.MustRegister(serverStreamMsgReceived)
	prom.MustRegister(serverStreamMsgSent)
	prom.MustRegister(serverHandledHistogram)
	prom.MustRegister(serverErred)
}

type ServerMonitor struct {
}

func (m *ServerMonitor) NewServerMonitor(rpcType monitoring.RpcType, fullMethod string) monitoring.RpcMonitor {
	r := &serverRpcMonitor{rpcType: rpcType, startTime: time.Now()}
	r.serviceName, r.methodName = splitMethodName(fullMethod)
	serverStartedCounter.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName).Inc()
	return r
}

type serverRpcMonitor struct {
	rpcType     monitoring.RpcType
	serviceName string
	methodName  string
	startTime   time.Time
}

func (r *serverRpcMonitor) ReceivedMessage() {
	serverStreamMsgReceived.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName).Inc()
}

func (r *serverRpcMonitor) SentMessage() {
	serverStreamMsgSent.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName).Inc()
}

func (r *serverRpcMonitor) Handled(code codes.Code) {
	serverHandledHistogram.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, code.String()).Observe(time.Since(r.startTime).Seconds())
}

func (r *serverRpcMonitor) Erred(err error) {
	serverErred.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, errorType(err)).Inc()
}

func splitMethodName(fullMethodName string) (string, string) {
	fullMethodName = strings.TrimPrefix(fullMethodName, "/") // remove leading slash
	if i := strings.Index(fullMethodName, "/"); i >= 0 {
		return fullMethodName[:i], fullMethodName[i+1:]
	}
	return "unknown", "unknown"
}

func errorType(err error) string {
	switch err.(type) {
	case transport.ConnectionError:
		return "ConnectionError"
	case transport.StreamError:
		return "StreamError"
	default:
		return "Unknown"
	}
}

package sdk

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func gRPCKeepAliveOpt() grpc.DialOption {
	return grpc.WithKeepaliveParams(keepalive.ClientParameters{
		// send pings every 10 seconds if there is no activity
		Time: 10 * time.Second,
		// wait 2 seconds for ping ack before considering the connection
		// dead
		Timeout: 2 * time.Second,
		// send pings even without active streams
		PermitWithoutStream: true,
	})
}

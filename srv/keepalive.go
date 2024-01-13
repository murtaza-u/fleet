package srv

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func gRPCKeepAliveOpts() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			// if a client pings more than once every 5 seconds,
			// terminate the connection
			MinTime: 5 * time.Second,
			// allow pings even when there are no active streams
			PermitWithoutStream: true,
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			// ping the client if it is idle for 5 seconds to ensure the
			// connection is still active
			Time: 5 * time.Second,
			// wait 2 seconds for the ping ack before assuming the
			// connection is dead
			Timeout: 2 * time.Second,
		}),
	}
}

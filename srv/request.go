package srv

import "github.com/murtaza-u/fleet/internal/pb"

// request wraps a protobuf request, holding the original (*pb.Request)
// and a reply channel for gRPC server to communicate the response to
// HTTP server.
type request struct {
	*pb.Request
	reply chan<- *pb.Response
}

package srv

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/murtaza-u/fleet/internal/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	ErrSubdomainInUse   = status.Error(codes.AlreadyExists, "subdomain already in use")
	ErrSubdomainMissing = status.Error(codes.InvalidArgument, "missing subdomain")
	ErrSubdomainInvalid = status.Error(codes.InvalidArgument, "invalid subdomain")
)

// Listen implements the gRPC Fleet service's Listen method. It
// facilitates communication between HTTP requests and connected gRPC
// clients, acting as a relay in between.
func (s Srv) Listen(stream pb.Fleet_ListenServer) error {
	reply, err := stream.Recv()
	if err != nil {
		return status.Error(codes.Unknown, err.Error())
	}
	subdomain := reply.GetSubdomain()
	if err := s.verifySubdomain(subdomain); err != nil {
		return err
	}
	err = stream.Send(&pb.Call{
		Payload: &pb.Call_ServingAt{
			ServingAt: fmt.Sprintf(s.servingUrlFormat, subdomain),
		},
	})
	if err != nil {
		return status.Error(codes.Unknown, err.Error())
	}

	queue := make(chan request)
	defer close(queue)

	s.store.Put(subdomain, queue)
	defer s.store.Delete(subdomain)

	errCh := make(chan error, 1)
	defer close(errCh)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	defer wg.Wait()

	go func() {
		defer wg.Done()
		for {
			err := s.route(stream)
			if err != nil {
				errCh <- err
				return
			}
		}
	}()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case err := <-errCh:
			return err
		case req := <-queue:
			s.store.Put(req.Request.GetId(), req.reply)
			err := stream.Send(&pb.Call{
				Payload: &pb.Call_Request{
					Request: req.Request,
				},
			})
			if err != nil {
				return status.Error(codes.Unknown, err.Error())
			}
		}
	}
}

// route receives a response from the connected gRPC client and routes
// it to the intended invoker(HTTP client).
func (s Srv) route(stream pb.Fleet_ListenServer) error {
	reply, err := stream.Recv()
	if err != nil {
		return status.Error(codes.Unknown, err.Error())
	}
	res := reply.GetResponse()
	if res == nil {
		return status.Error(codes.InvalidArgument, "empty response")
	}
	id := res.GetId()
	if id == "" {
		return status.Error(codes.InvalidArgument,
			"missing client id in response")
	}
	v := s.store.Get(id)
	if v == nil {
		return status.Error(codes.Internal, "an internal error occured")
	}
	tunnel, ok := v.(chan<- *pb.Response)
	if !ok {
		return status.Error(codes.Internal, "an internal error occured")
	}
	tunnel <- res
	s.store.Delete(id)
	return nil
}

func (s Srv) verifySubdomain(subdomain string) error {
	if subdomain == "" {
		return ErrSubdomainMissing
	}
	if strings.Contains(subdomain, ".") {
		return ErrSubdomainInvalid
	}
	if s.store.Exists(subdomain) {
		return ErrSubdomainInUse
	}
	return nil
}

// runRPC starts the gRPC server, blocking until an error occurs.
func (s Srv) runRPC() error {
	ln, err := net.Listen("tcp", s.rpcPort)
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	if s.tls {
		certs, err := loadTLSCreds(s.cert, s.priv)
		if err != nil {
			return err
		}
		opts = []grpc.ServerOption{grpc.Creds(certs)}
	}
	opts = append(opts, gRPCKeepAliveOpts()...)
	opts = append(opts,
		grpc.MaxRecvMsgSize(s.maxRecvSize),
		grpc.MaxSendMsgSize(s.maxSendSize))

	grpcS := grpc.NewServer(opts...)
	pb.RegisterFleetServer(grpcS, s)
	if s.reflect {
		reflection.Register(grpcS)
	}
	return grpcS.Serve(ln)
}

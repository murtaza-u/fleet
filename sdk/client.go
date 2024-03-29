package sdk

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/murtaza-u/fleet/internal/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// client is the Fleet gRPC client implementation.
type client struct {
	opts
	stream pb.Fleet_ListenClient
}

// newClient creates a new client.
func newClient(ctx context.Context, optfns ...OptFunc) (*client, error) {
	o := defaultOpts()
	for _, fn := range optfns {
		fn(&o)
	}
	if err := o.validate(); err != nil {
		return nil, err
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if o.tls {
		creds, err := loadTLSCreds(o.caCert)
		if err != nil {
			return nil, err
		}
		opts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	}
	opts = append(opts, gRPCKeepAliveOpt())
	c, err := grpc.Dial(o.rpcAddr, opts...)
	if err != nil {
		return nil, err
	}

	stream, err := pb.NewFleetClient(c).Listen(ctx)
	if err != nil {
		if stat, ok := status.FromError(err); ok {
			return nil, fmt.Errorf("[%s] %s", stat.Code(), stat.Message())
		}
		return nil, err
	}

	return &client{
		opts:   o,
		stream: stream,
	}, nil
}

// ListenAndServe serves requests, blocking until an error occurs or the
// connection is interrupted.
func (c client) ListenAndServe(h http.Handler) error {
	defer c.close()

	wg := new(sync.WaitGroup)
	defer wg.Wait()

	if err := c.registerSubdomain(); err != nil {
		return err
	}

	for {
		payload, err := c.stream.Recv()
		if err != nil {
			if stat, ok := status.FromError(err); ok {
				return fmt.Errorf("[%s] %s", stat.Code(), stat.Message())
			}
			return err
		}
		req := payload.GetRequest()
		if req == nil {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			w, err := c.process(h, req)
			if err != nil {
				return
			}
			err = c.stream.Send(&pb.Reply{
				Payload: &pb.Reply_Response{
					Response: &pb.Response{
						Id:     req.GetId(),
						Data:   w.buffer.Bytes(),
						Status: int32(w.statusCode),
						Header: w.httpToPbHeader(),
					},
				},
			})
			if err != nil {
				return
			}
		}()
	}
}

func (c client) registerSubdomain() error {
	err := c.stream.Send(&pb.Reply{
		Payload: &pb.Reply_Subdomain{
			Subdomain: c.subdomain,
		},
	})
	if err != nil {
		if stat, ok := status.FromError(err); ok {
			return fmt.Errorf("[%s] %s", stat.Code(), stat.Message())
		}
		return err
	}
	payload, err := c.stream.Recv()
	if err != nil {
		if stat, ok := status.FromError(err); ok {
			return fmt.Errorf("[%s] %s", stat.Code(), stat.Message())
		}
		return err
	}
	servingAt := payload.GetServingAt()
	if servingAt == "" {
		return fmt.Errorf("unexpected payload received")
	}
	log.Printf("serving at %q\n", servingAt)
	return nil
}

func (client) process(h http.Handler, req *pb.Request) (*rw, error) {
	body := bytes.NewReader(req.GetBody())
	r, err := http.NewRequest(req.GetMethod(), req.GetUrl(), body)
	if err != nil {
		return nil, err
	}
	for k, values := range req.GetHeader() {
		for _, v := range values.Fields {
			r.Header.Add(k, v)
		}
	}
	w := newResponseWriter()
	h.ServeHTTP(w, r)
	return w, nil
}

func (c client) close() error { return c.stream.CloseSend() }

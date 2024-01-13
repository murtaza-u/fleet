package srv

import (
	"github.com/murtaza-u/fleet/internal/pb"

	"github.com/murtaza-u/dream"
)

// Srv is the Fleet gRPC server implementation.
type Srv struct {
	opts
	store *dream.Store

	pb.UnimplementedFleetServer
}

// New creates a new instance of Srv.
func New(optfns ...OptFunc) (*Srv, error) {
	o := defaultOpts()
	for _, fn := range optfns {
		err := fn(&o)
		if err != nil {
			return nil, err
		}
	}
	return &Srv{
		opts:  o,
		store: dream.New(),
	}, nil
}

// Run starts the gRPC and HTTP servers, blocking until an error occurs.
// It returns an error if either server fails.
func (s Srv) Run() error {
	err := make(chan error)
	go func() { err <- s.runRPC() }()
	go func() { err <- s.runHTTP() }()
	return <-err
}

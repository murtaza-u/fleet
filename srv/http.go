package srv

import (
	"io"
	"net/http"

	"github.com/murtaza-u/fleet/internal/pb"
)

// runHTTP starts the HTTP server, blocking until an error occurs.
func (s Srv) runHTTP() error {
	http.HandleFunc("/", s.httpHandler)
	return http.ListenAndServe(s.httpPort, nil)
}

// httpHandler proxies incoming requests to the connected gRPC client.
func (s Srv) httpHandler(w http.ResponseWriter, r *http.Request) {
	subdomain := subdomain(r.Host)
	if subdomain == "" {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}
	v := s.store.Get(subdomain)
	if v == nil {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}
	queue, ok := v.(chan request)
	if !ok {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, s.maxRequestBodySize))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	r.Body.Close()
	header := make(map[string]*pb.HeaderFields, len(r.Header))
	for k, v := range r.Header {
		header[k] = &pb.HeaderFields{Fields: v}
	}
	req := &pb.Request{
		Id:     newID(),
		Method: r.Method,
		Url:    r.URL.String(),
		Body:   body,
		Header: header,
	}
	reply := make(chan *pb.Response, 1)
	queue <- request{
		Request: req,
		reply:   reply,
	}
	res := <-reply
	close(reply)
	for k, values := range res.GetHeader() {
		for _, v := range values.GetFields() {
			w.Header().Add(k, v)
		}
	}
	if res.GetStatus() != 0 {
		w.WriteHeader(int(res.GetStatus()))
	}
	w.Write(res.GetData())
}

package main

import (
	"fmt"
	"log"
	"net/http"

	fleet "github.com/murtaza-u/fleet/sdk"

	"github.com/go-chi/chi/v5"
)

const preferredSubdomain = "foo"

func main() {
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		v := q.Get("foo")
		w.WriteHeader(http.StatusOK)
		addTxt := "\nYou can also try hitting `/foo`\n"
		if v == "" {
			w.Write([]byte("index" + addTxt))
			return
		}
		w.Write([]byte("index|" + "foo=" + v + addTxt))
	})
	router.Get("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("foo\n"))
	})

	fmt.Println("------")
	fmt.Printf("Attempting to connect to %q via gRPC\n", "localhost:2035")
	fmt.Printf("If the connection is established, try hitting http://%s.%s\n",
		preferredSubdomain, "localhost:8080")
	fmt.Println("------")

	log.Fatal(fleet.Handle(router,
		fleet.WithRPCAddress("localhost:2035"),
		fleet.WithPreferredSubdomain(preferredSubdomain),
	))
}

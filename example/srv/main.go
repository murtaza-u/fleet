package main

import (
	"fmt"
	"log"

	"github.com/murtaza-u/fleet/srv"
)

func main() {
	srv, err := srv.New(srv.WithReflection())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("------")
	fmt.Printf("gRPC address: %s\n", "localhost:2035")
	fmt.Printf("http address: %s\n", "localhost:8080")
	fmt.Printf("tls: %s\n", "disabled")
	fmt.Println("------")
	fmt.Println("Run the client in another terminal window")

	log.Fatal(srv.Run())
}

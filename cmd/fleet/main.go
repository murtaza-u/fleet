package main

import (
	"log"

	"github.com/murtaza-u/fleet/cli"
)

func init() {
	log.SetFlags(0)
}

func main() {
	err := cli.Run()
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"github.com/alwindoss/casa/internal/server"
)

func main() {
	log.Fatal(server.Run())
}

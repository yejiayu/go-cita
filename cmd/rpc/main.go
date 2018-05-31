package main

import (
	"log"

	"github.com/yejiayu/go-cita/rpc"
)

func main() {
	if err := rpc.New(":8080"); err != nil {
		log.Fatal(err)
	}
}

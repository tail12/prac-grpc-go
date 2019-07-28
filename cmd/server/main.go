package main

import (
	"fmt"
	"github.com/tail12/prac-grpc-go/pkg/server"
	"os"
)

func main() {
	if err := server.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

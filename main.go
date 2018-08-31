package main

import (
	"fmt"
	"os"

	"github.com/ykore52/rest_reversi/server"
)

func main() {

	fmt.Println("Startup server.")

	if err := server.Run(8080, os.Args); err != nil {
		fmt.Printf("Server failed: %s", err.Error())
	}
}

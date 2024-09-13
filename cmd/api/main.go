package main

import (
	"BuntServer/internal/server"
	"fmt"
)

func main() {
	server := server.NewServer()

	err := server.ListenAndServe()

	if err != nil {
		panic(fmt.Sprintf("Error starting server: %s", err))
	}
}

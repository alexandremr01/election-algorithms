package main

import (
	"fmt"
	"os"
)

func main() {
	nodeID := os.Getenv("NODE_ID")
    fmt.Printf("My ID: %s\n", nodeID)
}

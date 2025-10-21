package main

import (
	"fmt"
)

func main() {
	fmt.Println("ROUTING MICROSERVICE")

	// TODO: Run the KAFKA consumer for it to wait for upcoming messages

	// TODO: Send the output produced by RoutingService as a new msg on another KAFKA topic, acting as producer
}
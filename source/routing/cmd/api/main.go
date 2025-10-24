package main

import (
	"context"
	"fmt"
	"geopathplanner/routing/internal/kafka"
	"geopathplanner/routing/internal/service"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

// This is the main function to be called
// It will reads for new messages on RequestTopic, processes them with the function and sending the results on ResponseTopic
func Run() error {
	// Read from environment variables
	brokersEnv := os.Getenv("KAFKA_BROKERS")
	requestTopic := os.Getenv("KAFKA_REQUEST_TOPIC")
	responseTopic := os.Getenv("KAFKA_RESPONSE_TOPIC")
	groupID := os.Getenv("KAFKA_GROUP_ID")

	fmt.Println("Kafka configuration:")
	fmt.Printf("  KAFKA_BROKERS        = %s\n", brokersEnv)
	fmt.Printf("  KAFKA_REQUEST_TOPIC  = %s\n", requestTopic)
	fmt.Printf("  KAFKA_RESPONSE_TOPIC = %s\n", responseTopic)
	fmt.Printf("  KAFKA_GROUP_ID       = %s\n", groupID)

	if brokersEnv == "" || requestTopic == "" || responseTopic == "" || groupID == "" {
		return fmt.Errorf("missing required Kafka environment variables")
	}

	// Split broker string into slice
	brokers := strings.Split(brokersEnv, ",")

	// Create KafkaService
	k, err := kafka.NewKafkaService(context.Background(), brokers, groupID, requestTopic, responseTopic)
	if err != nil {
		return fmt.Errorf("error while creating KafkaService: %v", err)
	}

	// Create RoutingService
	_, err = service.NewRoutingService()
	if err != nil {
		return fmt.Errorf("error while creating RoutingService: %v", err)
	}

	// ========== MAIN LOOP =============
	// TODO: Run the KAFKA consumer for it to wait for upcoming messages
	k.ConsumeMessage(func(r *kgo.Record) {
		// Generate random number between 1 to 10
		n := rand.Intn(10) + 1
		log.Printf("Handling request (%ds)...\n", n)
		// TODO: Just for testing, simulate Ns cooldown and then send a response
		time.Sleep(time.Duration(n) * time.Second)

		k.ProduceMessage([]byte(fmt.Sprintf(
			"Response to topic=%s partition=%d offset=%d key=%s value=%s\n",
			r.Topic,
			r.Partition,
			r.Offset,
			string(r.Key),
			string(r.Value),
		)))

		// TODO: Implement validation
		// // 1. Validate data to make sure it is a valide RoutingRequest
		// req, err := models.NewRoutingRequestFromJson(string(r.Value))
		// if err != nil {
		// 	fmt.Printf("Error decoding RoutingRequest (topic=%s, partition=%d, offset=%d): %v",
		// 		r.Topic, r.Partition, r.Offset, err)
		// 	return
		// }

		// fmt.Printf("Received RoutingRequest %s with %d waypoints",
		// 	req.RequestID, len(req.Waypoints))

		// // 2. Run RoutingService
		// response := rs.HandleRoutingRequest(req)
		
		// // 3. Marshal response to obtain []byte
		// data, err := json.Marshal(response)
		// if err != nil {
		// 	fmt.Printf("Failed to marshal RoutingResponse: %v", err)
		// 	return
		// }

		// // 4. Send the output produced by RoutingService as a new msg on another KAFKA topic, acting as producer
		// k.ProduceMessage(data)
	})	

	return nil
}

func main() {
	fmt.Printf("ROUTING MICROSERVICE\n\n")

	if err := Run(); err != nil {
		fmt.Printf("abort: %v\n", err)
		os.Exit(1)
	}
}
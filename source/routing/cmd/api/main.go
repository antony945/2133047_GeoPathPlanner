package main

import (
	"context"
	"encoding/json"
	"fmt"
	"geopathplanner/routing/internal/kafka"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/service"
	"geopathplanner/routing/internal/validator"
	"os"
	"strings"

	"github.com/twmb/franz-go/pkg/kgo"
)

// This is the main function to be called
// It will reads for new messages on RequestTopic, processes them with the function and sending the results on ResponseTopic
func Run() error {
	// Read from environment variables
	brokersEnv := os.Getenv("KAFKA_BROKERS")
	requestTopic := os.Getenv("KAFKA_REQUEST_TOPIC")
	responseTopic := os.Getenv("KAFKA_RESPONSE_TOPIC")
	groupID := os.Getenv("KAFKA_ROUTING_CONSUMER_GROUP_ID")

	fmt.Println("Kafka configuration:")
	fmt.Printf("  KAFKA_BROKERS        = %s\n", brokersEnv)
	fmt.Printf("  KAFKA_REQUEST_TOPIC  = %s\n", requestTopic)
	fmt.Printf("  KAFKA_RESPONSE_TOPIC = %s\n", responseTopic)
	fmt.Printf("  KAFKA_ROUTING_CONSUMER_GROUP_ID       = %s\n", groupID)

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
	defer k.Close()

	// Create RoutingService
	rs, err := service.NewRoutingService()
	if err != nil {
		return fmt.Errorf("error while creating RoutingService: %v", err)
	}

	// ========== MAIN LOOP =============
	// Run the KAFKA consumer for it to wait for upcoming messages
	k.ConsumeMessage(func(r *kgo.Record) error {
		// What to do with each record that we have
		// Extract request_id from header "request_id"
		// 📨 Log the basic info
	    fmt.Printf("📩 Received message from topic %s (partition=%d, offset=%d)\n",
        r.Topic, r.Partition, r.Offset)

		// ✅ Extract request_id from headers
		var requestID string
		for _, h := range r.Headers {
			if h.Key == "request_id" {
				requestID = string(h.Value)
				break
			}
		}

		if requestID == "" {
			fmt.Println("⚠️ No request_id header found — skipping message")
			return nil // or handle according to your logic
		}

		fmt.Printf("🆔 Extracted request_id: %s\n", requestID)

		var response *models.RoutingResponse

		// 1. Validate data to make sure it is a valide RoutingRequest
		v := validator.NewDefaultValidator()
		req, err := v.ValidateMessage(r.Value)
		if err != nil {
			error_msg := fmt.Sprintf("❌ error decoding RoutingRequest (topic=%s, partition=%d, offset=%d): %v",
				r.Topic, r.Partition, r.Offset, err)
			fmt.Println(error_msg)
			// TODO: What to do if we get something that is not a routingRequest?
			// Alternative: send an error msg and continue
			// TODO: For now just return a notfound route
			// take received at form
			
			response = models.NewRoutingResponseError(models.MustNewEmptyRoutingRequest(requestID, r.Timestamp), err.Error())
		} else {
			fmt.Printf("✅ Valid RoutingRequest %s with %d wps and %d constraints (topic=%s, partition=%d, offset=%d)\n", req.RequestID, len(req.Waypoints), len(req.Constraints), r.Topic, r.Partition, r.Offset)

			// 2. Run RoutingService
			myResp, found := rs.HandleRoutingRequest(req, v)
			if found {
				fmt.Printf("✅ ROUTE FOUND\n")
			} else {
				fmt.Printf("❌ ROUTE NOT FOUND\n")
			}

			response = myResp
		}

		// 3. Marshal response to obtain []byte
		data, err := json.Marshal(response)
		if err != nil {
			fmt.Printf("failed to marshal RoutingResponse: %v", err)
			return nil
		}

		// 4. Send the output produced by RoutingService as a new msg on another KAFKA topic, acting as producer
		k.ProduceMessage(data)
		return nil
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
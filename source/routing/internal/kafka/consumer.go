package kafka

// Consume message from a specific topic and pass it to the RoutingService to handle it.
// Return the output of the RoutingService to be later used in the producer.
// Implement consumer in a way that a worker pool can exist, each processing one or more messages at the "same time"
package kafka

// Given the output of the RoutingService, act as a producer by sending it onto a kafka topic so for the backend to read it.
// Do it in a coordinate way as the consumer (maybe same number of workers in working pool?)
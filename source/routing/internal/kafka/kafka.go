package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// Consume message from a specific topic and pass it to the RoutingService to handle it.
// Return the output of the RoutingService to be later used in the producer.
// Implement consumer in a way that a worker pool can exist, each processing one or more messages at the "same time"

// Given the output of the RoutingService, act as a producer by sending it onto a kafka topic so for the backend to read it.
// Do it in a coordinate way as the consumer (maybe same number of workers in working pool?)

func TestKafka() error {
	topic := "test-go-topic"
    partition := 0
    conn, err := connect(topic, partition)
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}

    writeMessages(conn, []string{"msg 1", "msg 22","msg 333"})
    readMessages(conn, 10, 10e3)
    readWithReader(topic, "consumer-through-kafka 1")

    if err := conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
    }

    return nil
}

//Connect to the specified topic and partition in the server
// Non-containerized apps can connect through: localhost:9092
// Dockerized apps con connect through: kafka:9093
func connect(topic string, partition int)(*kafka.Conn, error){
    conn, err := kafka.DialLeader(context.Background(), "tcp", 
        "kafka:9093", topic, partition)
    if err != nil {
        fmt.Println("failed to dial leader")
    }
    return conn, err
} //end connect

//Writes the messages in the string slice to the topic
func writeMessages(conn *kafka.Conn, msgs []string){
    var err error
    conn.SetWriteDeadline(time.Now().Add(10*time.Second)) 

    for _, msg := range msgs{
        _, err = conn.WriteMessages(
        kafka.Message{Value: []byte(msg)},)
    }
    if err != nil {
        fmt.Println("failed to write messages:", err)
    }
} //end writeMessages

//Reads all messages in the partition from the start
//Specify a minimum and maximum size in bytes to read (1 char = 1 byte)
func readMessages(conn *kafka.Conn, minSize int, maxSize int){
    conn.SetReadDeadline(time.Now().Add(5*time.Second))
    batch := conn.ReadBatch(minSize, maxSize) //in bytes

    msg:= make([]byte, 10e3)     //set the max length of each message
    for {
        msgSize, err := batch.Read(msg)
        if err != nil {
            break
        }
        fmt.Printf("CONSUMED MESSAGE -> %s\n", string(msg[:msgSize]))
    }

    if err := batch.Close(); err != nil {   //make sure to close the batch
        fmt.Println("failed to close batch:", err)
    }
} //end readMessages

//Read from the topic using kafka.Reader
//Readers can use consumer groups (but are not required to)
func readWithReader(topic string, groupID string){
    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers:    []string{"localhost:9092", "localhost:9093", "localhost:9094"},
        // Brokers:    []string{"localhost:9092", "kafka:9093"},
        GroupID:    groupID,
        Topic:      topic,
        MaxBytes:   100,     //per message
        // more options are available
    })

    //Create a deadline
    readDeadline, _ := context.WithDeadline(context.Background(),
        time.Now().Add(5*time.Second))
    for {
        msg, err := r.ReadMessage(readDeadline)
        if err != nil {
            break
        }
        fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", 
            msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
    }

    if err := r.Close(); err != nil {
        fmt.Println("failed to close reader:", err)
    }
}
package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/twmb/franz-go/pkg/kgo"
)

type KakfaService struct {
	Client        *kgo.Client
	RequestTopic  string
	ResponseTopic string
	GroupID       string
	Brokers       []string
	Ctx context.Context
}

func NewKafkaService(ctx context.Context, brokers []string, groupID, requestTopic, responseTopic string) (*KakfaService, error) {

	// Create the client
	// One client can both produce and consume!
	// Consuming can either be direct (no consumer group), or through a group. Below, we use a group.
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(requestTopic),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka client: %w", err)
	}

	fmt.Printf("Successfully created Kafka client with '%s' groupID listening on '%s' topic..\n", groupID, requestTopic)

	return &KakfaService{
		Client: client,
		RequestTopic: requestTopic,
		ResponseTopic: responseTopic,
		GroupID: groupID,
		Brokers: brokers,
		Ctx: ctx,
	}, nil
}

func (k *KakfaService) ProduceMessage(data []byte) error {
	fmt.Printf("Producing msg...\n\n")

	// ctx := context.Background()

	// 1.) Producing a message
	// All record production goes through Produce, and the callback can be used
	// to allow for synchronous or asynchronous production.
	var wg sync.WaitGroup
	wg.Add(1)

	// Create record to put on topic
	record := &kgo.Record{Topic: k.ResponseTopic, Value: data}
	
	k.Client.Produce(k.Ctx, record, func(_ *kgo.Record, err error) {
		defer wg.Done()
		if err != nil {
			fmt.Printf("record had a produce error: %v\n", err)
		}

	})
	wg.Wait()

	// // Alternatively, ProduceSync exists to synchronously produce a batch of records.
	// if err := k.Client.ProduceSync(ctx, record).FirstErr(); err != nil {
	// 	fmt.Printf("record had a produce error while synchronously producing: %v\n", err)
	// }

	return nil
}

func (k *KakfaService) ConsumeMessage(handleRecord func(*kgo.Record) error) {
	fmt.Printf("Waiting for messages to be consumed...\n\n")
	
	// 2.) Consuming messages from a topic
	for {
		fetches := k.Client.PollFetches(k.Ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			// All errors are retried internally when fetching, but non-retriable errors are
			// returned from polls so that users can notice and take action.
			panic(fmt.Sprint(errs))
		}

		// or a callback function.
		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			for _, record := range p.Records {
				// What to do with each record that we have
				fmt.Printf(
					"> Consumed from topic=%s partition=%d offset=%d key=%s value=%s\n",
					record.Topic,
					record.Partition,
					record.Offset,
					string(record.Key),
					string(record.Value),
				)

				err := handleRecord(record)
				if err != nil {
					// TODO: Stop listening
				}
			}
		})
	}
}

func (k *KakfaService) Close() error {
	k.Client.Close()
	return nil
}

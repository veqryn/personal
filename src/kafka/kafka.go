package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
)

func main() {

	config := sarama.NewConfig()

	brokers := []string{"kafka:9092"}

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = producer.Close(); err != nil {
			panic(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	var enqueued, errors int
	doneCh := make(chan struct{})
	go func() {
		for {

			time.Sleep(500 * time.Millisecond)

			strTime := strconv.Itoa(int(time.Now().Unix()))
			msg := &sarama.ProducerMessage{
				Topic: "important",
				Key:   sarama.StringEncoder(strTime),
				Value: sarama.StringEncoder("Something Cool"),
			}
			select {
			case producer.Input() <- msg:
				enqueued++
				fmt.Println("Produce message")
			case producerErr := <-producer.Errors():
				errors++
				fmt.Println("Failed to produce message:", producerErr)
			case <-signals:
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	log.Printf("Enqueued: %d; errors: %d\n", enqueued, errors)
}

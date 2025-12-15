package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	cfg "github.com/pobyzaarif/go-config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	RabbitMQURL string `env:"RABBITMQ_URL"`
}

type ProductMessage struct {
	ProductCode string `json:"product_code"`
	ProductName string `json:"product_name"`
	Stock       int    `json:"stock"`
}

func main() {
	config := Config{}
	cfg.LoadConfig(&config)
	queue := "product-queue"

	conn, err := amqp.Dial(config.RabbitMQURL)
	if err != nil {
		log.Fatalf("dial err: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel err: %v", err)
	}
	defer ch.Close()

	// declare queue
	_, err = ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("queue declare err: %v", err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < 10; i++ {
		wg.Add(1)
		i := i // capture loop var
		go func() {
			defer wg.Done()

			// message JSON
			msg := ProductMessage{
				ProductCode: fmt.Sprintf("12345000%d", i),
				ProductName: "soy sauce",
				Stock:       12,
			}

			payload, err := json.Marshal(msg)
			if err != nil {
				log.Printf("json marshal err: %v", err)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			mu.Lock()
			err = ch.PublishWithContext(ctx,
				"",    // default exchange
				queue, // routing key = queue
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        payload,
				},
			)
			mu.Unlock()

			if err != nil {
				log.Printf("publish err: %v", err)
				return
			}

			log.Println("Message published:", string(payload))
		}()
	}

	wg.Wait()
}

package main

import (
	"context"
	"log"
	"os"

	banye "github.com/Egot3/Banye"
	"github.com/Egot3/Dialyn/internal/manager"
	"github.com/Egot3/Dialyn/internal/middleware"
	diacon "github.com/Egot3/Zhao"
	"github.com/Egot3/Zhao/queues"
	"github.com/Egot3/Zhao/sub"
)

func main() {
	config := diacon.RabbitMQConfiguration{
		URL:  os.Getenv("RABBIT_URL"),
		Port: os.Getenv("RABBIT_PORT"),
	}
	conn, err := diacon.Connect(config)
	if err != nil {
		log.Panicf("Couldn't connect: %v", err)
	}
	defer conn.Close()

	subscriber, err := sub.NewSubscriber(conn)
	if err != nil {
		log.Panicf("Couldn't create a subscriber: %v", err)
	}
	defer subscriber.Close()

	subscriberManager := manager.NewSubscriberManager(subscriber)

	updateChan := make(chan []*queues.QueueStruct)

	client := banye.NewClient(nil)
	before, after := middleware.TraceTripperMiddleware()
	client.UseTripper(before, after)

	go manager.CheckForUpdates(context.Background(), client, updateChan)
	go subscriberManager.Reconcile(updateChan)
}

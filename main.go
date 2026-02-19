package main

import (
	"log"
	"net/http"

	wss "github.com/Egot3/Dialyn/internal/wssConnection"
	diacon "github.com/Egot3/Zhao"
	"github.com/Egot3/Zhao/queues"
	"github.com/Egot3/Zhao/sub"
)

func main() {
	config := diacon.RabbitMQConfiguration{
		URL:  "amqp://guest:guest@localhost",
		Port: "1130",
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

	qStruct := queues.QueueStruct{
		Name:           "RINGRINGRING",
		Durable:        false,
		DeleteOnUnused: false,
		Exclusive:      false,
		NoWait:         false,
		Args:           nil,
	}
	q, err := queues.NewQueue(subscriber.Ch, qStruct)
	if err != nil {
		log.Panicf("Couldn't create a queue: %v", err)
	}

	f, err := subscriber.StartSubscriberFunc(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Panicf("Couldn't get a starter func: %v", err)
	}
	forever := make(chan any, 1)

	go f(forever)

	http.HandleFunc("/wss", wss.WssHandler(forever))
	log.Fatal(http.ListenAndServe(":8250", nil))

	// for message := range forever {
	// 	log.Printf("got message: %#v", message)
	// }

	<-forever
}

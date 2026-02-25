package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	wss "github.com/Egot3/Dialyn/internal/wssConnection"
	diacon "github.com/Egot3/Zhao"
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

	f, err := subscriber.StartSubscriberFunc("test-queue", "", true, false, false, false, nil)
	if err != nil {
		log.Panicf("Couldn't get a starter func: %v", err)
	}
	forever := make(chan any, 1)

	go f(forever)

	http.HandleFunc("/wss", wss.WssHandler(forever))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("OWN_PORT")), nil))

	// for message := range forever {
	// 	log.Printf("got message: %#v", message)
	// }

	<-forever
}

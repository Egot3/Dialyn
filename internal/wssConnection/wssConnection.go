package wss

import (
	"context"
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	amqp "github.com/rabbitmq/amqp091-go"
)

func WssHandler(ch chan any) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})
		if err != nil {
			log.Panicf("couldn't upgrade connection")
		}
		defer conn.CloseNow()

		log.Println("CONNECTED VIA WSS")

		ctx := conn.CloseRead(context.Background())

		log.Println("Sending messages!")
		for v := range ch {
			log.Printf("Message: %v", v)
			value := v.(amqp.Delivery)
			log.Printf("Received: %#v", value)

			if ctx.Err() != nil {
				log.Println("connection closed, stopping")
				return
			}

			if err := wsjson.Write(ctx, conn, string(value.Body)); err != nil {
				log.Printf("write error: %v", err)
				return
			}
		}

		conn.Close(1000, "Doesn't want to continue")
	}
}

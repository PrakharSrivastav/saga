package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"syreclabs.com/go/faker"
)

const (
	PaymentChannel    string = "PaymentChannel"
	OrderChannel      string = "OrderChannel"
	DeliveryChannel   string = "DeliveryChannel"
	RestaurantChannel string = "RestaurantChannel"
	ReplyChannel      string = "ReplyChannel"

	StartMsg    string = "Start"
	DoneMsg     string = "DoneMsg"
	ErrorMsg    string = "ErrorMsg"
	RollbackMsg string = "CreateMsg"
)

func main() {
	var err error
	mux := http.NewServeMux()
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	if _, err = client.Ping().Result(); err != nil {
		log.Fatalf("error creating redis client %s", err)
	}

	o := &Orchestrator{
		c: client,
		r: client.Subscribe(PaymentChannel, OrderChannel, DeliveryChannel, RestaurantChannel, ReplyChannel),
	}
	go o.start()

	mux.HandleFunc("/create", o.create)
	log.Println("starting server")

	log.Fatal(http.ListenAndServe(":8080", mux))
}

type Orchestrator struct {
	c *redis.Client
	r *redis.PubSub
}

func (o Orchestrator) create(writer http.ResponseWriter, request *http.Request) {
	if _, err := fmt.Fprintf(writer, "responding"); err != nil {
		log.Printf("error while writing %s", err.Error())
	}
	o.send()
}

func (o Orchestrator) send() {
	transactionID := faker.Bitcoin().Address()
	if err := o.c.Publish(OrderChannel, struct {
		TransactionId string `json:"transaction_id"`
		MsgType       string `json:"msg_type"`
	}{
		TransactionId: transactionID,
		MsgType:       StartMsg,
	}).Err(); err != nil {
		log.Printf("error publishing to the channel %s\n", err)
	}
}

func (o Orchestrator) start() {
	if _, err := o.r.Receive(); err != nil {
		log.Fatalf("error setting up redis %s \n", err)
	}

	ch := o.r.Channel()
	log.Println("starting the redis client")
	for {
		select {
		case msg := <-ch:
			log.Printf("message is %v \n", msg)
		}
	}
}

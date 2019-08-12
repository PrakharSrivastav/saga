package main

import (
	"encoding/json"
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

	ServicePayment    string = "Payment"
	ServiceOrder      string = "Order"
	ServiceRestaurant string = "Restaurant"
	ServiceDelivery   string = "Delivery"

	ActionStart    string = "Start"
	ActionDone     string = "DoneMsg"
	ActionError    string = "ErrorMsg"
	ActionRollback string = "RollbackMsg"
)

type Message struct {
	ID      string `json:"id"`
	Service string `json:"service"`
	Action  string `json:"action"`
	Message string `json:"message"`
}

func (m Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

type Orchestrator struct {
	c *redis.Client
	r *redis.PubSub
}

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

func (o Orchestrator) create(writer http.ResponseWriter, request *http.Request) {
	if _, err := fmt.Fprintf(writer, "responding"); err != nil {
		log.Printf("error while writing %s", err.Error())
	}
	m := Message{
		ID:      faker.Bitcoin().Address(),
		Message: "Something",
	}
	o.next(OrderChannel, ServiceOrder, m)
}

func (o Orchestrator) start() {
	if _, err := o.r.Receive(); err != nil {
		log.Fatalf("error setting up redis %s \n", err)
	}

	ch := o.r.Channel()
	log.Println("starting the redis client")
	defer o.r.Close()
	for {
		select {
		case msg := <-ch:
			m := Message{}
			err := json.Unmarshal([]byte(msg.Payload), &m)
			if err != nil {
				log.Println(err)
				continue
			}

			// only process the messages on ReplyChannel
			switch msg.Channel {
			case ReplyChannel:
				// if there is any error, just rollback
				if m.Action != ActionDone {
					log.Printf("Rolling back transaction with id %s", m.ID)
					o.rollback(m)
					continue
				}

				// else start the next stage
				switch m.Service {
				case ServiceOrder:
					o.next(PaymentChannel, ServicePayment, m)
				case ServicePayment:
					o.next(RestaurantChannel, ServiceRestaurant, m)
				case ServiceRestaurant:
					o.next(DeliveryChannel, ServiceDelivery, m)
				case ServiceDelivery:
					log.Println("Food Delivered")
				}
			}
		}
	}
}

func (o Orchestrator) next(channel, service string, message Message) {
	var err error
	message.Action = ActionStart
	message.Service = service
	if err = o.c.Publish(channel, message).Err(); err != nil {
		log.Printf("error publishing start-message to %s channel", channel)
	}
	log.Printf("start message published to channel :%s", channel)
}

func (o Orchestrator) rollback(m Message) {
	var err error
	message := Message{
		ID:      m.ID,
		Action:  ActionRollback,
		Message: "May day !! May day!!",
	}
	if err = o.c.Publish(OrderChannel, message).Err(); err != nil {
		log.Printf("error publishing rollback message to %s channel", OrderChannel)
	}
	if err = o.c.Publish(PaymentChannel, message).Err(); err != nil {
		log.Printf("error publishing rollback message to %s channel", PaymentChannel)
	}
	if err = o.c.Publish(RestaurantChannel, message).Err(); err != nil {
		log.Printf("error publishing rollback message to %s channel", RestaurantChannel)
	}
	if err = o.c.Publish(DeliveryChannel, message).Err(); err != nil {
		log.Printf("error publishing rollback message to %s channel", DeliveryChannel)
	}

}

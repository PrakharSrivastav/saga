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
	ActionStart       string = "Start"
	ActionDone        string = "DoneMsg"
	ActionError       string = "ErrorMsg"
	ActionRollback    string = "RollbackMsg"
)

// Message represents the payload sent over redis pub/sub
type Message struct {
	ID      string `json:"id"`
	Service string `json:"service"`
	Action  string `json:"action"`
	Message string `json:"message"`
}

// MarshalBinary should be implemented to send message to redis
func (m Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// Orchestrator orchestrates the order processing task
type Orchestrator struct {
	c *redis.Client
	r *redis.PubSub
}

func main() {
	var err error
	// create client and ping redis
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	if _, err = client.Ping().Result(); err != nil {
		log.Fatalf("error creating redis client %s", err)
	}

	// initialize and start the orchestrator in the background
	o := &Orchestrator{
		c: client,
		r: client.Subscribe(PaymentChannel, OrderChannel, DeliveryChannel, RestaurantChannel, ReplyChannel),
	}
	go o.start()

	// define ServeMux and start the server
	mux := http.NewServeMux()
	mux.HandleFunc("/create", o.create)
	log.Println("starting server")

	log.Fatal(http.ListenAndServe(":8080", mux))
}

// create the order - step 1 in the order processing pipeline
func (o Orchestrator) create(writer http.ResponseWriter, request *http.Request) {
	if _, err := fmt.Fprintf(writer, "responding"); err != nil {
		log.Printf("error while writing %s", err.Error())
	}
	m := Message{ID: faker.Bitcoin().Address(), Message: "Something"}
	o.next(OrderChannel, ServiceOrder, m)
}

// start the goroutine to orchestrate the process
func (o Orchestrator) start() {
	var err error
	if _, err = o.r.Receive(); err != nil {
		log.Fatalf("error setting up redis %s \n", err)
	}
	ch := o.r.Channel()
	defer func() { _ = o.r.Close() }()

	log.Println("starting the redis client")
	for {
		select {
		case msg := <-ch:
			m := Message{}
			if err = json.Unmarshal([]byte(msg.Payload), &m); err != nil {
				log.Println(err)
				// continue to skip bad messages
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

// next triggers start operation on the other micro-services
// based on the channel and service
func (o Orchestrator) next(channel, service string, message Message) {
	var err error
	message.Action = ActionStart
	message.Service = service
	if err = o.c.Publish(channel, message).Err(); err != nil {
		log.Printf("error publishing start-message to %s channel", channel)
	}
	log.Printf("start message published to channel :%s", channel)
}

// rollback instructs the other micro-services to rollback the transaction
func (o Orchestrator) rollback(m Message) {
	var err error
	message := Message{
		ID:      m.ID, // ID is mandatory
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

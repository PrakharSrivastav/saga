package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis"
	"syreclabs.com/go/faker"
)

const (
	OrderChannel   string = "OrderChannel"
	ReplyChannel   string = "ReplyChannel"
	ActionStart    string = "Start"
	ActionDone     string = "DoneMsg"
	ActionRollback string = "RollbackMsg"
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

func main() {
	// create client and ping redis
	var err error
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	if _, err = client.Ping().Result(); err != nil {
		log.Fatalf("error creating redis client %s", err)
	}

	// subscribe to the required channels
	pubsub := client.Subscribe(OrderChannel, ReplyChannel)
	if _, err = pubsub.Receive(); err != nil {
		log.Fatalf("error subscribing %s", err)
	}
	defer pubsub.Close()
	ch := pubsub.Channel()

	log.Println("starting the order service")
	for {
		select {
		case msg := <-ch:
			m := Message{}
			err := json.Unmarshal([]byte(msg.Payload), &m)
			if err != nil {
				log.Println(err)
				continue
			}

			switch msg.Channel {
			case OrderChannel:
				log.Printf("recieved message with id %s ", m.ID)

				// random sleep to simulate some work in action
				d := faker.RandomInt(1, 3)
				time.Sleep(time.Duration(d) * time.Second)

				// Happy Flow
				if m.Action == ActionStart {
					m.Action = ActionDone
					if err = client.Publish(ReplyChannel, m).Err(); err != nil {
						log.Printf("error publishing done-message to %s channel", ReplyChannel)
					}
					log.Printf("done message published to channel :%s", ReplyChannel)
				}

				// Rollback flow
				if m.Action == ActionRollback {
					log.Printf("rolling back transaction with ID :%s", m.ID)
				}

			}
		}
	}

}

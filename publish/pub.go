package main

import (
	"log"

	"github.com/go-redis/redis"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("error connecting to redis %s", err)
	}
	log.Printf("response is %s \n", pong)

	pubsub := client.Subscribe("ch1", "ch2","PaymentChannel")
	_, err = pubsub.Receive()
	if err != nil {
		log.Fatalf("error subscribing %s", err)
	}

	ch := pubsub.Channel()

	for {
		select {
		case msg := <-ch:
			log.Println(msg.Channel)
			log.Println(msg.Payload)
			log.Println(msg.Pattern)
			log.Println("---------------")
		}
	}

}

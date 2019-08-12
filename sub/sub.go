package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"syreclabs.com/go/faker"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("error connecting to redis %s", err)
	}
	log.Printf("response is %s \n", pong)

	pubsub := client.Subscribe("ch1")
	_, err = pubsub.Receive()
	if err != nil {
		log.Fatalf("error subscribing %s", err)
	}

	for i := 0; i < 10; i++ {
		err = client.Publish("PaymentChannel", NewUser().String()).Err()
		err = client.Publish("ch2", NewUser().String()).Err()
		if err != nil {
			panic(err)
		}
	}

	time.AfterFunc(time.Second, func() {
		_ = pubsub.Close()
		log.Println("closing after a second")
	})

}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u User) String() string {
	return fmt.Sprintf("{ID:%s,Name:%s,Email:%s}", u.ID, u.Name, u.Email)
}

func NewUser() User {
	return User{
		ID:    faker.Bitcoin().Address(),
		Name:  faker.Internet().UserName(),
		Email: faker.Internet().Email(),
	}
}

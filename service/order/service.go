package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	var err error
	s := Service()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer log.Println("responded")
		if err = s.Action(); err != nil {
			log.Printf("err %#v", err)
		}
		if _, err = fmt.Fprintf(w, "order received"); err != nil {
			log.Printf("err %#v", err)
		}
	})
	log.Println("starting server")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

type Impl struct {
	conn string
}

func Service() contract.Service {
	return Impl{}
}

func (i Impl) Rollback() error {
	log.Println("Rolling back")
	return nil
}

func (i Impl) Propagate() error {
	log.Println("Propagate")
	return nil
}

func (i Impl) Action() error {
	log.Println("Acting")
	return nil
}

func (i Impl) Feedback() error {
	log.Println("Acting")
	return nil
}

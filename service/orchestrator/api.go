package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	var err error
	a := Service()
	if err = a.Listen(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer log.Println("responded")
		if _, err = fmt.Fprintf(w, "input"); err != nil {
			log.Printf("error getting order %#v\n", err)
		}
	})

	log.Println("starting server")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

type service interface {
	Listen() error
	Rollback() error
	Propagate() error
}

type Impl struct {
	conn string
}

func (Impl) Listen() error {
	log.Println("Listen")
	return nil
}

func (Impl) Rollback() error {
	log.Println("Rollback")
	return nil
}

func (Impl) Propagate() error {
	log.Println("Propagate")
	return nil
}

func Service() service {
	log.Println("initing service")
	return Impl{}
}

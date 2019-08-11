package main

import (
	"net/http"

	//"github.com/PrakharSrivastav/saga/stub/orchestrator/server"
	//sc "github.com/PrakharSrivastav/saga/stub/service/client"
	"github.com/labstack/gommon/log"
)

// order orchestrator
func main() {
	//addr := "localhost:8080"
	//clientHTTP := http.Client{
	//	Timeout:   2 * time.Second,
	//	Transport: &http.Transport{IdleConnTimeout: 5 * time.Second},
	//}

	log.Info("starting order orchestrator")
	//if err := server.Run(addr, OrderOrchestrator{client: clientHTTP}); err != nil {
	//	log.Fatal(err)
	//}
}

type OrderOrchestrator struct {
	client http.Client
}
//
//func (o OrderOrchestrator) Feedback(context.Context, *server.FeedbackRequest) (*server.Response, error) {
//	panic("implement me")
//	return nil, nil
//}
//
//func (o OrderOrchestrator) Next(context.Context, *server.NextRequest) (*server.Response, error) {
//	panic("implement me")
//	return nil, nil
//}
//
//func (o OrderOrchestrator) Rollback(context.Context, *server.RollbackRequest) (*server.Response, error) {
//	panic("implement me")
//	return nil, nil
//}
//
//func (o OrderOrchestrator) Start(context.Context, *server.StartRequest) (*server.Response, error) {
//	log.Info("Starting the Order Orchestrator")
//	c := sc.NewServiceClient("http://localhost:8081", &o.client)
//	ctx := context.Background()
//	_, err := c.Execute(ctx, nil)
//	if err != nil {
//		panic(err)
//	}
//	panic("yo")
//	return nil, nil
//}

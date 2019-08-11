package contract

type Orchestrator interface {
	Start(StartRequest) Response
	Feedback(FeedbackRequest) Response
	Next(NextRequest) Response
	Rollback(RollbackRequest) Response
}

type StartRequest struct {
	Type string `json:"type"`
}

type FeedbackRequest struct {
	ID      string `json:"id"`
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type NextRequest struct {
	ID string `json:"id"`
}

type RollbackRequest struct {
	ID string `json:"id"`
}

type Response struct {
	ID      string `json:"id"`
	IsError bool   `json:"is_error"`
	Message string `json:"message"`
}

/*
remoto generate  ./orchestrator/contract/contract.remoto.go ./orchestrator/templates/client.go.plush -o ./orchestrator/stubs/client.go \
&& gofmt -w ./orchestrator/stubs/client.go

remoto generate  ./orchestrator/contract/contract.remoto.go ./orchestrator/templates/server.go.plush -o ./orchestrator/stubs/server.go \
&& gofmt -w ./orchestrator/stubs/server.go
*/

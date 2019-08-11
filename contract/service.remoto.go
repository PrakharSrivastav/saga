package contract

type Service interface {
	Execute(request MSExecuteRequest) MSResponse
	Rollback(request MSRollbackRequest) MSResponse
}

type MSExecuteRequest struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type MSResponse struct {
	ID      string `json:"id"`
	IsError bool   `json:"is_error"`
	Message string `json:"message"`
}

type MSRollbackRequest struct {
	ID string `json:"id"`
}

/*
server
remoto generate  contract/service.remoto.go templates/server.go.plush -o stub/service/server/server.go \
&& gofmt -w stub/service/server/server.go


client

remoto generate  contract/service.remoto.go templates/client.go.plush -o stub/service/client/client.go \
&& gofmt -w stub/service/client/client.go
*/

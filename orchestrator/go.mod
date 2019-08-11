module github.com/PrakharSrivastav/saga/orchestrator

require (
	github.com/PrakharSrivastav/saga/stub/orchestrator/server v0.0.0
	github.com/PrakharSrivastav/saga/stub/service/client v0.0.0
	github.com/labstack/gommon v0.2.9
)

replace (
	github.com/PrakharSrivastav/saga/stub/orchestrator/server v0.0.0 => ../stub/orchestrator/server
	github.com/PrakharSrivastav/saga/stub/service/client v0.0.0 => ../stub/service/client
)

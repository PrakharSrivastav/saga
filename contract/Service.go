package contract

type Service interface {
	Action() error
	Feedback() error
	Propagate() error
	Rollback() error
}

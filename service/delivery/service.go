package delivery

import "log"

type Service interface {
	Action() error
	Rollback() error
	Propagate() error
}

type Impl struct {
	conn string
}

func New() Service {
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

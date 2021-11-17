package executor

import "context"

type Dummy struct {
	Name   string
	Result int
}

func (d *Dummy) Identifier() string {
	return d.Name
}

func (d *Dummy) Execute(_ context.Context) (int, error) {
	return d.Result, nil
}

func (d *Dummy) CacheEligible() bool {
	return false
}

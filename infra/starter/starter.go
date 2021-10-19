package starter

import (
	"context"

	"github.com/hashicorp/go-multierror"
)

// Starter is anything capable of being started. Starters, on this context, are expected to return
// only if something went wrong or the provided context has been cancelled.
type Starter interface {
	Start(context.Context) error
}

// Engine holds a series of starters and monitor them once they are started.
type Engine struct {
	starters []Starter
}

// NewEngine returns a new starter engine. This starter engine is capable of starting all provided
// Starters at once.
func NewEngine(starters ...Starter) *Engine {
	return &Engine{
		starters: starters,
	}
}

// Start starts all inner starters. If one of returns after a Start call all of them are cancelled
// as Starters are never expected to return (exception made for when they do error out or the
// context has been cancelled).
func (e *Engine) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error, len(e.starters))
	for _, s := range e.starters {
		go func(s Starter) {
			err := s.Start(ctx)
			cancel()
			errs <- err
		}(s)
	}

	var errors *multierror.Error
	for i := 0; i < len(e.starters); i++ {
		if err := <-errs; err != nil {
			errors = multierror.Append(errors, err)
		}
	}
	return errors.ErrorOrNil()
}

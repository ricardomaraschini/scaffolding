package starter

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type starter struct {
	misbehave bool
}

func (s starter) Start(ctx context.Context) error {
	if s.misbehave {
		return fmt.Errorf("i reckon i better not to behave")
	}
	<-ctx.Done()
	return nil
}

func TestStart(t *testing.T) {
	engine := NewEngine(
		starter{},
		starter{},
		starter{},
	)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	errs := make(chan error)
	go func() {
		errs <- engine.Start(ctx)
	}()

	select {
	case err := <-errs:
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	case <-time.NewTicker(10 * time.Second).C:
		t.Fatalf("timeout waiting for starter to end")
	}
}

func TestFailToStart(t *testing.T) {
	engine := NewEngine(
		starter{},
		starter{},
		starter{misbehave: true},
	)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	errs := make(chan error)
	go func() {
		errs <- engine.Start(ctx)
	}()

	select {
	case err := <-errs:
		if err == nil {
			t.Errorf("expected error to be returned, nil received instead")
		}
	case <-time.NewTicker(10 * time.Second).C:
		t.Fatalf("timeout waiting for starter to end")
	}
}

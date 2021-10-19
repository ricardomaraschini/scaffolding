package ctrls

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PromCtrlOpt sets an option in a PromCtrl.
type PromCtrlOpt func(*PromCtrl)

// WithPromCtrlBindAddress changes the default bind address for a PromCtrl object.
func WithPromCtrlBindAddress(bind string) PromCtrlOpt {
	return func(p *PromCtrl) {
		p.bind = bind
	}
}

// PromCtrl handles prometheus metric requests.
type PromCtrl struct {
	bind string
}

// NewPromCtrl returns a controller for metrics endpoint.
func NewPromCtrl(opts ...PromCtrlOpt) *PromCtrl {
	ctrl := &PromCtrl{
		bind: ":8181",
	}
	for _, opt := range opts {
		opt(ctrl)
	}
	return ctrl
}

// Start puts the http server online.
func (p *PromCtrl) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    p.bind,
		Handler: promhttp.Handler(),
	}

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("error shutting down https server: %s", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

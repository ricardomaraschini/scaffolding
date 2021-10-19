package ctrls

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// AppCtrlOpt sets an option in an AppCtrl.
type AppCtrlOpt func(*AppCtrl)

// WithAppCtrlBindAddress changes the default bind address for an AppCtrl object.
func WithAppCtrlBindAddress(bind string) AppCtrlOpt {
	return func(a *AppCtrl) {
		a.bind = bind
	}
}

// AppCtrl is our sample controller.
type AppCtrl struct {
	http.Handler
	bind string
}

// NewAppCtrl returns a controller for http requests.
func NewAppCtrl(opts ...AppCtrlOpt) *AppCtrl {
	ctrl := &AppCtrl{
		bind: ":8080",
	}

	router := mux.NewRouter()
	router.HandleFunc(
		"/", instrumentHandler("root", ctrl.Root),
	).Methods(http.MethodGet)
	router.HandleFunc(
		"/hello/{name}", instrumentHandler("hello", ctrl.Hello),
	)
	ctrl.Handler = router

	for _, opt := range opts {
		opt(ctrl)
	}
	return ctrl
}

// Hello handles requests to "/hello" endpoint.
func (a *AppCtrl) Hello(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello %q\n", vars["name"])
}

// Root handles requests to "/" endpoint.
func (a *AppCtrl) Root(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Start puts the http server online.
func (a *AppCtrl) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    a.bind,
		Handler: a,
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

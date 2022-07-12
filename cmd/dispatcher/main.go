package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

func main() {
	log.Info("dispatcher starting...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", healthz)
	mux.HandleFunc("/api/v1/multi/batches", dispatch)
	srv := http.Server{Addr: ":8000", Handler: mux}
	ctx := context.Background()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Info("shutting down dispatcher service...")
			srv.Shutdown(ctx)
			<-ctx.Done()
		}
	}()

	go moniter(ctx)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %s", err)
	}
}

//call by k8s liveness probe
func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	data := (time.Now()).String()
	log.Debug("health ok")
	w.Write([]byte(data))
}

func dispatch(w http.ResponseWriter, req *http.Request) {

	go execute(req.Context())
}

func execute(ctx context.Context) {

}

func moniter(ctx context.Context) {

}

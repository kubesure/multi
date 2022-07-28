package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kubesure/multi"
	"github.com/kubesure/multi/internal"
	"github.com/kubesure/multi/internal/dispatcher"
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
	mux.HandleFunc("/api/v1/multi/batches/jobs", dispatch)
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
	j, err := parseJob(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		erres := multi.Erroresponse{Code: err.Code, Message: err.Message}
		data, _ := json.Marshal(erres)
		fmt.Fprintf(w, "%s", data)
	} else {
		saverr := dispatcher.SaveJob(j)
		if saverr != nil {
			erres := multi.Erroresponse{Code: saverr.Code, Message: saverr.Message}
			data, _ := json.Marshal(erres)
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "%s", data)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func execute(ctx context.Context) {

}

func moniter(ctx context.Context) {

}

func parseJob(req *http.Request) (*internal.Job, *multi.Error) {
	body, _ := ioutil.ReadAll(req.Body)
	j := internal.Job{}
	err := json.Unmarshal([]byte(body), &j)
	if err != nil {
		log.Errorf("err %v during unmarshalling data %s ", err, body)
		return nil, &multi.Error{Code: multi.HTTPError, Message: multi.HTTPRequestError}
	}
	return &j, nil
}

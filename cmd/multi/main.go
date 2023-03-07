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
	internal "github.com/kubesure/multi/internal"
	mi "github.com/kubesure/multi/internal/multi"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

func main() {

	r := mux.NewRouter()
	//r.Headers("content-Type", "application/json")
	r.HandleFunc("/", healthz).Methods("GET")
	r.HandleFunc("/api/v1/multi/batches", scheduleBatch).Methods("POST")
	r.HandleFunc("/api/v1/multi/batches/{id}", scheduledBatchInfo).Methods("GET")
	r.HandleFunc("/api/v1/multi/batches/{id}", updateBatchResult).Methods("PUT")
	r.HandleFunc("/api/v1/multi/batches/jobs/{id}", jobResult).Methods("GET")
	//r.MethodNotAllowedHandler = MethodNotAllowedHandler()
	http.Handle("/", r)

	srv := http.Server{
		Addr:         ":8000",
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	ctx := context.Background()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Info("shutting down multi dispatcher service...")
			srv.Shutdown(ctx)
			<-ctx.Done()
		}
	}()

	log.Info("multi service started...")

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

//Scheduler schedules requests on dispatchers
func scheduleBatch(w http.ResponseWriter, req *http.Request) {
	cs, err := parseCustomer(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		erres := multi.ErroResponse{Code: err.Code, Message: err.Message}
		data, _ := json.Marshal(erres)
		fmt.Fprintf(w, "%s", data)
	} else {
		id, err1 := mi.SaveBatch(internal.CustomerSearchType, cs)
		if err1 != nil {
			erres := multi.ErroResponse{Code: err1.Code, Message: err1.Message}
			data, _ := json.Marshal(erres)
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "%s", data)
		}
		w.WriteHeader(http.StatusCreated)
		// write to location header
		w.Write([]byte("batch " + id))
	}

}

func scheduledBatchInfo(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	batch, err := mi.GetBatch(vars["id"])
	if err == nil {
		erres := multi.ErroResponse{Code: err.Code, Message: err.Message}
		data, _ := json.Marshal(erres)
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "%s", data)
	} else {
		data, _ := json.Marshal(batch)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", data)
	}
}

func updateBatchResult(w http.ResponseWriter, req *http.Request) {
	j, err := praseJob(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		erres := multi.ErroResponse{Code: err.Code, Message: err.Message}
		data, _ := json.Marshal(erres)
		fmt.Fprintf(w, "%s", data)
	} else {
		err1 := mi.UpdateJob(j)
		if err1 != nil {
			erres := multi.ErroResponse{Code: err1.Code, Message: err1.Message}
			data, _ := json.Marshal(erres)
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "%s", data)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func jobResult(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	data := (time.Now()).String() + "result"
	log.Debug("result")
	w.Write([]byte(data))
}

func parseCustomer(req *http.Request) (*internal.CustomerSearch, *multi.Error) {
	body, _ := ioutil.ReadAll(req.Body)
	cs := internal.CustomerSearch{}
	err := json.Unmarshal([]byte(body), &cs)
	if err != nil {
		log.Errorf("err %v during unmarshalling data %s ", err, body)
		return nil, &multi.Error{Code: multi.HTTPError, Message: multi.HTTPRequestError}
	}
	return &cs, nil
}

func praseJob(req *http.Request) (*internal.Job, *multi.Error) {
	body, _ := ioutil.ReadAll(req.Body)
	j := internal.Job{}
	err := json.Unmarshal([]byte(body), &j)
	if err != nil {
		log.Errorf("err %v during unmarshalling data %s ", err, body)
		return nil, &multi.Error{Code: multi.HTTPError, Message: multi.HTTPRequestError}
	}
	return &j, nil
}

func MethodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Method not allowed")
	})
}

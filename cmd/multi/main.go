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

	internal "github.com/kubesure/multi/internal"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type CustomerScheduleRes struct {
}

func main() {

	r := mux.NewRouter()
	//r.Headers("content-Type", "application/json")
	r.HandleFunc("/", healthz).Methods("GET")
	r.HandleFunc("/api/v1/multi/searches/customers", scheduleCustomerSearch).Methods("POST")
	r.HandleFunc("/api/v1/multi/searches/customers/{id}", scheduledBatchInfo).Methods("GET")
	r.HandleFunc("/api/v1/multi/searches/customers/{id}/searches/{id}", updateSearchResult).Methods("PUT")
	r.MethodNotAllowedHandler = MethodNotAllowedHandler()
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

//Scheduler schedules requests on dispatchers
func scheduleCustomerSearch(w http.ResponseWriter, req *http.Request) {
	cs, err := parseCustomer(req)
	if err != nil {
		log.Errorf("Error parsing customer")
	}
	id, err1 := internal.SaveBatch(internal.CustomerSearchType, cs)
	if err1 != nil {
		log.Errorf("Error saving customer")
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("scheduled " + id))
}

func scheduledBatchInfo(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("schedule returned for %v", vars["id"])))
}

func updateSearchResult(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	data := (time.Now()).String()
	log.Debug("health ok")
	w.Write([]byte(data))
}

func parseCustomer(req *http.Request) (internal.CustomerSearch, error) {
	body, _ := ioutil.ReadAll(req.Body)
	cs := internal.CustomerSearch{}
	err := json.Unmarshal([]byte(body), &cs)
	if err != nil {
		log.Errorf("err %v during unmarshalling data %s ", err, body)
	}
	return cs, err
}

//call by k8s liveness probe
func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	data := (time.Now()).String()
	log.Debug("health ok")
	w.Write([]byte(data))
}

func MethodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Method not allowed")
	})
}

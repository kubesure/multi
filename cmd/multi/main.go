package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kubesure/multi"
	"github.com/kubesure/multi/internal"
	log "github.com/sirupsen/logrus"
)

type createBatchReq struct {
	Jobs []internal.Job
}

type updateJobReq struct {
	Job internal.Job
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

func main() {

	router := gin.Default()
	router.GET("/healthz", healthz)
	router.POST("/api/v1/multi/batches", scheduleBatch)
	router.PUT("/api/v1/multi/batches/:id/jobs/:id", updateJob)
	router.GET("/api/v1/multi/batches/:id", scheduledBatchInfo)
	router.GET("/api/v1/multi/batches/:id/jobs/:id", jobInfo)

	srv := &http.Server{
		Addr:         ":8000",
		Handler:      router,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}

// call by k8s liveness probe
func healthz(c *gin.Context) {
	c.Writer.WriteHeader(200)
	data := (time.Now()).String()
	log.Debug("health ok")
	c.Writer.Write([]byte(data))
}

// Scheduler schedules requests on dispatchers
func scheduleBatch(c *gin.Context) {
	cbr, err := praseJobs(c.Request)
	if err != nil {
		erres := multi.ErroResponse{Code: err.Code, Message: err.Message}
		data, _ := json.Marshal(erres)
		c.String(http.StatusBadRequest, string(data))
	} else {
		batch, err1 := internal.SaveBatch(cbr.Jobs)

		if err1 != nil {
			erres := multi.ErroResponse{Code: err.Code, Message: err.Message}
			data, _ := json.Marshal(erres)
			c.String(http.StatusServiceUnavailable, string(data))
		} else {
			c.Writer.WriteHeader(http.StatusCreated)
			// write to location header
			data, err := response(batch)
			if err != nil {
				erres := multi.ErroResponse{Code: err.Code, Message: err.Message}
				data, _ := json.Marshal(erres)
				c.String(http.StatusServiceUnavailable, string(data))
			} else {
				c.String(http.StatusCreated, string(data))
			}
		}
	}
}

func scheduledBatchInfo(c *gin.Context) {
	id := c.Param("id")

	if len(id) == 0 {
		erres := multi.ErroResponse{Code: multi.HTTPError, Message: multi.HTTPRequestError}
		data, _ := json.Marshal(erres)
		c.String(http.StatusBadRequest, string(data))
	} else {
		batch, err1 := internal.GetBatch(id)
		if err1 != nil {
			erres := multi.ErroResponse{Code: err1.Code, Message: err1.Message}
			data, _ := json.Marshal(erres)
			c.String(http.StatusServiceUnavailable, string(data))
		}
		data, _ := json.Marshal(batch)
		c.String(http.StatusOK, string(data))
	}
}

func jobInfo(c *gin.Context) {

	var bid, jid string

	if len(c.Params) != 2 {
		erres := multi.ErroResponse{Code: multi.HTTPError, Message: multi.HTTPRequestError}
		data, _ := json.Marshal(erres)
		c.String(http.StatusBadRequest, string(data))
	} else {
		bid = c.Params[0].Value
		jid = c.Params[1].Value
		job, err1 := internal.GetJob(jid, bid)
		if err1 != nil {
			erres := multi.ErroResponse{Code: err1.Code, Message: err1.Message}
			data, _ := json.Marshal(erres)
			c.String(http.StatusServiceUnavailable, string(data))
		}

		data, err := internal.MarshalAny(job)
		if err != nil {
			erres := multi.ErroResponse{Code: err.Code, Message: err.Message}
			data, _ := json.Marshal(erres)
			c.String(http.StatusServiceUnavailable, string(data))
		} else {
			c.String(http.StatusOK, string(data))
		}
	}
}

func updateJob(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	job, err := internal.UnmarshalAny[internal.Job](body)

	if err != nil {
		erres := multi.ErroResponse{Code: err.Code, Message: err.Message}
		data, _ := json.Marshal(erres)
		c.String(http.StatusBadRequest, string(data))
	} else {
		job.Id = c.Params[1].Value
		err1 := internal.UpdateJob(job)
		if err1 != nil {
			erres := multi.ErroResponse{Code: err1.Code, Message: err1.Message}
			data, _ := json.Marshal(erres)
			c.String(http.StatusServiceUnavailable, string(data))
		} else {
			c.Writer.WriteHeader(http.StatusOK)
		}
	}
}

func response(batch *internal.Batch) (data []byte, err *multi.ErroResponse) {
	log := multi.NewLogger()
	data, errj := json.Marshal(batch)
	if errj != nil {
		log.LogInternalError(errj.Error())
		return nil, &multi.ErroResponse{Code: multi.HTTPError, Message: multi.HTTPResponseError}
	}
	return data, nil
}

func praseJobs(req *http.Request) (*createBatchReq, *multi.Error) {
	body, _ := ioutil.ReadAll(req.Body)
	var cbr createBatchReq
	err := json.Unmarshal([]byte(body), &cbr)
	if err != nil {
		log.Errorf("err %v during unmarshalling data %s ", err, body)
		return nil, &multi.Error{Code: multi.HTTPError, Message: multi.HTTPRequestError}
	}
	return &cbr, nil
}

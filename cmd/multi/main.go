package main

import (
	"context"
	"errors"
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

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

func main() {

	router := gin.Default()
	router.Use(internal.PreChecks(), internal.BeforeResponse())
	router.GET("/healthz", healthz)
	router.POST("/api/v1/multi/jobs", scheduleBatch)
	router.PUT("/api/v1/multi/batches/:id/jobs/:id", updateJob)
	router.GET("/api/v1/multi/batches/:id", scheduledBatchInfo)
	router.GET("/api/v1/multi/batches/:id/jobs/:id", jobInfo)

	srv := &http.Server{
		Addr:         ":8000",
		Handler:      router,
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
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

// Scheduler schedules requests on dispatchers
func scheduleBatch(c *gin.Context) {
	cbr := internal.CreateBatchReq{}
	if err := c.BindJSON(&cbr); err != nil {
		err := multi.ErrorResponse{Code: multi.HTTPError, Message: multi.InputInvalid}
		c.AbortWithStatusJSON(http.StatusBadRequest, internal.ResponseError(err, nil))
		return
	} else {
		batch, err1 := internal.SaveBatch(cbr.Jobs)
		if err1 != nil {
			erres := multi.ErrorResponse{Code: err1.Code, Message: err1.Message}
			c.JSON(http.StatusInternalServerError, erres)
		} else {
			// write to location header
			c.JSON(http.StatusCreated, batch)
		}
	}
}

func updateJob(c *gin.Context) {
	job := &internal.Job{}
	if err := c.BindJSON(job); err != nil {
		err := multi.ErrorResponse{Code: multi.HTTPError, Message: multi.InputInvalid}
		c.AbortWithStatusJSON(http.StatusBadRequest, internal.ResponseError(err, nil))
		return
	} else {
		job.Id = c.Params[1].Value
		err1 := internal.UpdateJob(job)
		if err1 != nil {
			erres := multi.ErrorResponse{Code: err1.Code, Message: err1.Message}
			c.JSON(http.StatusInternalServerError, erres)
		} else {
			c.Writer.WriteHeader(http.StatusOK)
		}
	}
}

func scheduledBatchInfo(c *gin.Context) {
	id := c.Param("id")
	if len(id) == 0 {
		err := multi.ErrorResponse{Code: multi.HTTPError, Message: multi.InputInvalid}
		c.AbortWithStatusJSON(http.StatusBadRequest, internal.ResponseError(err, nil))
		return
	} else {
		batch, err1 := internal.GetBatch(id)
		if err1 != nil {
			erres := multi.ErrorResponse{Code: err1.Code, Message: err1.Message}
			c.JSON(http.StatusInternalServerError, erres)
		} else {
			if batch == nil {
				err := multi.ErrorResponse{Code: multi.BatchNotFoundError, Message: multi.BatchNotFound}
				c.AbortWithStatusJSON(http.StatusNotFound, internal.ResponseError(err, nil))
			} else {
				c.JSON(http.StatusOK, batch)
			}
		}
	}
}

func jobInfo(c *gin.Context) {
	var bid, jid string
	if len(c.Params) != 2 {
		err := multi.ErrorResponse{Code: multi.HTTPError, Message: multi.InputInvalid}
		c.AbortWithStatusJSON(http.StatusBadRequest, internal.ResponseError(err, nil))
		return
	} else {
		bid = c.Params[0].Value
		jid = c.Params[1].Value
		job, err1 := internal.GetJob(jid, bid)
		if err1 != nil {
			erres := multi.ErrorResponse{Code: err1.Code, Message: err1.Message}
			c.JSON(http.StatusInternalServerError, erres)
		} else {
			if job == nil {
				err := multi.ErrorResponse{Code: multi.JobNotFoundError, Message: multi.JobNotFound}
				c.AbortWithStatusJSON(http.StatusNotFound, internal.ResponseError(err, nil))
			} else {
				c.JSON(http.StatusOK, job)
			}
		}
	}
}

// call by k8s liveness probe
func healthz(c *gin.Context) {
	c.Writer.WriteHeader(200)
}

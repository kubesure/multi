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

	router.POST("/api/v1/multi/batches/:id/jobs/", saveJob)

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
	log.Println("shutting down dispatcher server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("dispatcher server shutdown:", err)
	}
}

// Saves job in dispatcher db
func saveJob(c *gin.Context) {
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

// call by k8s liveness probe
func healthz(c *gin.Context) {
	c.Writer.WriteHeader(200)
}

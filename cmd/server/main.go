package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/tunedev/GoShortLink/pkg/handler"
	"github.com/tunedev/GoShortLink/pkg/store"
)

func main() {
	var PORT, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		PORT = 9921
	}

	mongoUri := os.Getenv("MONGO_URI")
	if mongoUri == "" {
		mongoUri = "mongodb://localhost:27017"
	}
	logger := log.New(os.Stdout, "web: ", log.LstdFlags)
	s, err := store.NewStore(mongoUri)
	if err != nil {
		logger.Fatalf("Failed to create a new Store: %v", err)
	}

	h := handler.NewHandler(s)

	readTimeout, err := strconv.Atoi(os.Getenv("READ_TIME_OUT"))
	if err != nil {
		readTimeout = 15
	}
	writeTimeout, err := strconv.Atoi(os.Getenv("WRITE_TIME_OUT"))
	if err != nil {
		writeTimeout = 15
	}
	srv := &http.Server {
		Handler: h,
		Addr: fmt.Sprintf(":%v", PORT),
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
	}
	
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Printf("Starting server on :%v", PORT)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("could not listen on :%v: %v\n", PORT, err)
		}
	}()

	// Block until a signal is received
	<- stop
	logger.Println("Shutting down the server...")

	shutDownDelay, err := strconv.Atoi(os.Getenv("SHUTDOWN_DELAY"))
	if err != nil {
		shutDownDelay = 15
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutDownDelay) * time.Second)
	defer cancel()
	srv.Shutdown(ctx)

	logger.Println("Server gracefully stopped")
}
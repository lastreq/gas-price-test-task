package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/lastreq/gas-price-test-task/internal/app/controlers"
	"github.com/lastreq/gas-price-test-task/internal/app/providers"
	"github.com/lastreq/gas-price-test-task/internal/app/services"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	prv := providers.New()

	svc := services.New(prv)

	ctr := controlers.New(svc)

	router := mux.NewRouter()
	router.HandleFunc("/get-gas-info", ctr.GetGasInfo).Methods("GET")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")
}

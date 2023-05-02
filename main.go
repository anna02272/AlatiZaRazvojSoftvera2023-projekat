package main

import (
	"context"
	"github.com/anna02272/AlatiZaRazvojSoftvera2023-projekat/config"
	service2 "github.com/anna02272/AlatiZaRazvojSoftvera2023-projekat/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var service *service2.Service

func main() {

	service := &service2.Service{
		Configurations: []*config.Config{},
	}
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)

	router.HandleFunc("/configurations", service.AddConfiguration).Methods("POST")
	router.HandleFunc("/configurations/{id}", service.GetConfiguration).Methods("GET")
	router.HandleFunc("/configurations/{id}", service.DeleteConfiguration).Methods("DELETE")
	// start server
	srv := &http.Server{Addr: "0.0.0.0:8000", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}

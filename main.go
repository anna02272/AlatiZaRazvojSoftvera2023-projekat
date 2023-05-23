// Config API
//
//	Title: Config API
//
//	Schemes: http
//	Version: 0.0.1
//	BasePath: /
//
//	Produces:
//	  - application/json
//
// swagger:meta
package main

import (
	"context"
	"github.com/anna02272/AlatiZaRazvojSoftvera2023-projekat/config"
	"github.com/anna02272/AlatiZaRazvojSoftvera2023-projekat/poststore"
	"github.com/anna02272/AlatiZaRazvojSoftvera2023-projekat/service"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	ps, err := poststore.New()
	if err != nil {
		log.Fatal(err)
	}

	service := &service.Service{
		Configurations: []*config.Config{},
		PostStore:      ps,
	}
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)

	router.HandleFunc("/configurations", service.AddConfiguration).Methods("POST")
	router.HandleFunc("/configurations/{id}/{version}", service.GetConfiguration).Methods("GET")
	router.HandleFunc("/configurations/{id}/{version}", service.DeleteConfiguration).Methods("DELETE")
	router.HandleFunc("/group", service.AddConfigurationGroup).Methods("POST")
	router.HandleFunc("/group/{id}/{version}", service.GetConfigurationGroup).Methods("GET")
	router.HandleFunc("/group/{id}/{version}", service.DeleteConfigurationGroup).Methods("DELETE")
	router.HandleFunc("/group/{id}/{version}/extend", service.ExtendConfigurationGroup).Methods("POST")
	router.HandleFunc("/swagger.yaml", service.SwaggerHandler).Methods("GET")

	// SwaggerUI
	optionsDevelopers := middleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	developerDocumentationHandler := middleware.SwaggerUI(optionsDevelopers, nil)
	router.Handle("/docs", developerDocumentationHandler)

	// ReDoc
	// optionsShared := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	// sharedDocumentationHandler := middleware.Redoc(optionsShared, nil)
	// router.Handle("/docs", sharedDocumentationHandler)

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

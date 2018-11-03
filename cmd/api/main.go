package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paulbes/go-pedal/pedal/client"

	"github.com/paulbes/go-pedal/pedal"
	api "github.com/paulbes/go-pedal/pkg/api/server"
	store "github.com/paulbes/go-pedal/pkg/api/store/http"
	"github.com/paulbes/go-pedal/pkg/md"
)

var clientIdentifier string

func init() {
	flag.StringVar(&clientIdentifier, "client-identifier", "", "Oslo City Bike Client Identifier")
	flag.Parse()
}

func main() {
	// Create an HTTP client for interacting with the city bike API
	cli, err := client.NewHTTPClient(clientIdentifier, 5)
	if err != nil {
		log.Fatalf("failed to create an API client: %s", err)
	}
	pedlar := pedal.New(cli)

	// Create a store that uses the pedlar interface
	stationStore := store.NewStationStore(pedlar)

	// Create a service that reads from the store
	stationService := api.NewStationService(stationStore)

	// Create the endpoints that interact with the known services
	services := api.Services{
		Station: stationService,
	}
	endpoints := api.MakeEndpoints(services)

	// Create HTTP handlers and attach them to routes so they
	// can be queried
	handlers := api.AttachRoutes(api.MakeHandlers(endpoints))

	// Create an entry point
	router := http.NewServeMux()
	// Ensure that cross origin requests are accepted
	http.Handle("/", md.Cors(router))
	// Add the known routes to the primary router
	router.Handle("/v1/", handlers)

	// Create an HTTP server
	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 300 * time.Second,
		Addr:         fmt.Sprintf(":%d", 8080),
	}

	// Start serving incoming HTTP requests
	errs := make(chan error, 2)
	go func() {
		log.Printf("server is listening on port %d", 8080)
		errs <- server.ListenAndServe()
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	log.Println("terminated", <-errs)
}

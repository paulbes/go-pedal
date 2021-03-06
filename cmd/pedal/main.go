package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/paulbes/go-pedal/pedal/client"

	"github.com/fatih/color"
	"github.com/paulbes/go-pedal/pedal"
)

var clientIdentifier string

func init() {
	flag.StringVar(&clientIdentifier, "client-identifier", "", "Oslo City Bike Client Identifier")
	flag.Parse()
}

func main() {
	// Create an http API client
	cli, err := client.NewHTTPClient(clientIdentifier, 5)
	if err != nil {
		log.Fatalf("failed to create an API client: %s", err)
	}

	// Read all stations
	stations, err := pedal.New(cli).Stations()
	if err != nil {
		log.Fatalf("failed to get stations: %s", err)
	}

	// Pretty print stations
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)
	w.Write([]byte(fmt.Sprintf("%-30s%-10s%-10s%s\n", "Station", "Locks", "Bikes", "Closed")))
	for _, station := range stations {
		out := []byte(fmt.Sprintf("%-30s%-10s%-10s%s\n",
			station.Title,
			color.CyanString(strconv.Itoa(station.Availability.Locks)),
			color.GreenString(strconv.Itoa(station.Availability.Bikes)),
			color.RedString(fmt.Sprintf("%t", station.Closed)),
		))
		_, err := w.Write(out)
		if err != nil {
			log.Fatalf("failed to write station information")
		}
	}
	err = w.Flush()
	if err != nil {
		log.Fatalf("failed to flush stations: %s", err)
	}
}

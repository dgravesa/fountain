package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/dgravesa/fountain/pkg/data"
	"github.com/dgravesa/fountain/pkg/handlers"
)

func main() {
	port := flag.Uint("port", 8080, "port to host server")

	reservoir := data.DefaultReservoir()
	waterlogs := handlers.NewWaterLogsHandler(reservoir)
	http.Handle("/waterlogs", waterlogs)

	portStr := fmt.Sprintf(":%d", *port)
	if err := http.ListenAndServe(portStr, nil); err != nil {
		log.Fatalln(err)
	}
}

package main

import (
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/lucabecci/stringsvc-microservices/internal"
	"github.com/lucabecci/stringsvc-microservices/services"
)

func main() {
	svc := services.GetService()
	upperCaseHandler := httptransport.NewServer(
		services.MakeUppercaseEndpoint(svc),
		services.DecodeUppercaseRequest,
		internal.EncodeResponse,
	)

	countHandler := httptransport.NewServer(
		services.MakeCountEndpoint(svc),
		services.DecodeCountRequest,
		internal.EncodeResponse,
	)

	http.Handle("/uppercase", upperCaseHandler)
	http.Handle("/count", countHandler)

	log.Fatal(http.ListenAndServe(":4000", nil))
}

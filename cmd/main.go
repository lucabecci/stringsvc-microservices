package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/lucabecci/stringsvc-microservices/internal"
	"github.com/lucabecci/stringsvc-microservices/services"
	"github.com/lucabecci/stringsvc-microservices/transports"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	var svc transports.StringService
	svc = transports.GetService()
	svc = internal.LoggingMiddleware{Logger: logger, Next: svc}
	upperCaseHandler := httptransport.NewServer(
		services.MakeUppercaseEndpoint(svc),
		transports.DecodeUppercaseRequest,
		internal.EncodeResponse,
	)
	countHandler := httptransport.NewServer(
		services.MakeCountEndpoint(svc),
		transports.DecodeCountRequest,
		internal.EncodeResponse,
	)

	http.Handle("/uppercase", upperCaseHandler)
	http.Handle("/count", countHandler)

	http.ListenAndServe(":4000", nil)
	fmt.Println("Server on port:4000")

}

func loggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			logger.Log("msg", "calling endpoint")
			defer logger.Log("msg", "called endpoint")
			return next(ctx, request)
		}
	}
}

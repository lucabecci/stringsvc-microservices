package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/lucabecci/stringsvc-microservices/internal"
	"github.com/lucabecci/stringsvc-microservices/services"
	"github.com/lucabecci/stringsvc-microservices/transports"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

	var svc transports.StringService
	svc = transports.GetService()
	svc = internal.LoggingMiddleware{Logger: logger, Next: svc}
	svc = internal.InstrumentingMiddleware{
		RequestCount:   requestCount,
		RequestLatency: requestLatency,
		CountResult:    countResult,
		Next:           svc,
	}

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
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", ":4000")
	logger.Log("err", http.ListenAndServe(":4000", nil))

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

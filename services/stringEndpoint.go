package services

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/lucabecci/stringsvc-microservices/transports"
)

func MakeUppercaseEndpoint(svc transports.StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(transports.UppercaseRequest)
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return transports.UppercaseResponse{V: v, Err: err.Error()}, nil
		}
		return transports.UppercaseResponse{V: v, Err: ""}, nil
	}
}

func MakeCountEndpoint(svc transports.StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(transports.CountRequest)
		v := svc.Count(req.S)
		return transports.CountResponse{V: v}, nil
	}
}

package services

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func MakeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UppercaseRequest)
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return UppercaseResponse{v, err.Error()}, nil
		}
		return UppercaseResponse{v, ""}, nil
	}
}

func MakeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CountRequest)
		v := svc.Count(req.S)
		return CountResponse{v}, nil
	}
}

package services

import (
	"errors"
	"strings"
)

type StringService interface {
	Uppercase(string) (string, error)
	Count(string) int
}

type stringService struct{}

func GetService() stringService {
	return stringService{}
}

func (stringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", errors.New("Send a word")
	}

	return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
	return len(s)
}

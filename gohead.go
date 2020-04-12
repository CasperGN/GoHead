package gohead

import (
	"net/http"
)

type GoHead struct{}

func NewGoHead() (*GoHead, error) {
	goHead := new(GoHead)

	return goHead, nil
}

func Probe(target string) (map[string][]string, string) {
	resp, _ := http.Get(target)

	var result map[string][]string
	for key, value := range resp.Header {
		result[key] = value
	}

	return result, target
}

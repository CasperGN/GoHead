package gohead

import (
	"crypto/tls"
	"net/http"
	"time"
)

type GoHead struct{}

func NewGoHead() (*GoHead, error) {
	goHead := new(GoHead)

	return goHead, nil
}

func Probe(target string) (map[string][]string, string) {
	client := &http.Client{
		Timeout: 8 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Proxy: http.ProxyFromEnvironment,
		},
	}

	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return nil, target
	}
	req.Header.Add("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		return nil, target
	}

	defer resp.Body.Close()

	return resp.Header, target
}

package gohead

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

type GoHead struct{}

func NewGoHead() (*GoHead, error) {
	goHead := new(GoHead)

	return goHead, nil
}

func Probe(target string) (io.ReadCloser, map[string][]string, string) {
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
		return nil, nil, target
	}
	req.Header.Add("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, target
	}

	defer resp.Body.Close()

	return resp.Body, resp.Header, target
}

package gohead

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

type GoHead struct{}

func NewGoHead() (*GoHead, error) {
	goHead := new(GoHead)

	return goHead, nil
}

func Probe(target string) /*(map[string][]string, string)*/ {
	fmt.Println(target)
	client := &http.Client{
		Timeout: 8,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Proxy: http.ProxyFromEnvironment,
		},
	}

	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		//return nil, target
	}
	req.Header.Add("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		//return nil, target
	}

	defer resp.Body.Close()

	fmt.Println(resp.Header)
	//return resp.Header, target
}

package gohead

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"
)

type GoHead struct{}

func NewGoHead() (*GoHead, error) {
	goHead := new(GoHead)

	return goHead, nil
}

func Probe(target string) (string, map[string][]string, string) {
	client := &http.Client{
		Timeout: 8 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // lgtm [go/disabled-certificate-check]
			Proxy:           http.ProxyFromEnvironment,
		},
	}

	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return "", nil, target
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:75.0) Gecko/20100101 Firefox/75.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", nil, target
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", resp.Header, target
	}

	return string(body), resp.Header, target
}

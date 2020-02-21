// Package http_helper contains helpers to interact with deployed resources through HTTP.
package http_helper

import (
	"crypto/tls"
	"fmt"
	"github.com/franroa/infratest/asserts"
	"github.com/franroa/infratest/retry"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
)

type HttpGet struct {
	url string
	t *testing.T
}

func (request HttpGet) Retry(retries int, sleepBetweenRetries time.Duration) asserts.AssertResponse {
	return HttpGetWithRetry(request.t, request.url, nil, retries, sleepBetweenRetries)
}

func HttpGetE(t *testing.T, url string, tlsConfig *tls.Config) (int, []byte, *http.Response, error) {
	logger.Logf(t, "Making an HTTP GET call to URL %s", url)

	// Set HTTP client transport config
	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := http.Client{
		// By default, Go does not impose a timeout, so an HTTP connection attempt can hang for a LONG time.
		Timeout: 10 * time.Second,
		// Include the previously created transport config
		Transport: tr,
	}

	resp, err := client.Get(url)
	if err != nil {
		return -1, nil, nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return -1, nil, nil, err
	}

	return resp.StatusCode, body, resp, nil
}



// HttpGetWithRetry repeatedly performs an HTTP GET on the given URL until the given status code and body are returned or until max
// retries has been exceeded.
func HttpGetWithRetry(t *testing.T, url string, tlsConfig *tls.Config, retries int, sleepBetweenRetries time.Duration) asserts.AssertResponse {
	response, err := HttpGetWithRetryE(t, url, tlsConfig, retries, sleepBetweenRetries)

	if err != nil {
		t.Fatal(err)
	}

	return response
}

// HttpGetWithRetryE repeatedly performs an HTTP GET on the given URL until the given status code and body are returned or until max
// retries has been exceeded.
func HttpGetWithRetryE(t *testing.T, url string, tlsConfig *tls.Config, retries int, sleepBetweenRetries time.Duration) (asserts.AssertResponse, error) {
	response, err := retry.DoWithRetryE(t, fmt.Sprintf("HTTP GET to URL %s", url), retries, sleepBetweenRetries, func() (interface{}, error) {
		status, body, resp, err := HttpGetE(t, url, tlsConfig)
		return asserts.GetResponse(body, status, resp, t), err
	})

	return response.(asserts.AssertResponse), err
}

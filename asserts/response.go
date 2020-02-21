package asserts

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

type AssertResponse struct {
	body         []byte
	bodyAsString string
	status       int
	response     *http.Response
	t            *testing.T
}

func GetResponse(body []byte, status int, response *http.Response, t *testing.T) AssertResponse {
	return AssertResponse{
		body:         body,
		bodyAsString: "",
		status:       status,
		response:     response,
		t:            t,
	}
}

func (response AssertResponse) HasStatus(status int) {
	assert.Equal(response.t, status, response.response.StatusCode)
}

func (response AssertResponse) Contains(body string) {
	if response.bodyAsString == "" {
		response.bodyAsString = strings.TrimSpace(string(response.body))
	}

	assert.Equal(response.t, body, response.bodyAsString)
}

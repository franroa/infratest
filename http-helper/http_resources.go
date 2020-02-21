package http_helper

import (
	"fmt"
	"github.com/franroa/infratest/k8s"
	"testing"
)

type KubernetesResourcesRequest struct {
	t *testing.T
	options *k8s.KubectlOptions
}

func Get(t *testing.T, options *k8s.KubectlOptions) KubernetesResourcesRequest {
	return KubernetesResourcesRequest{ t: t, options:options }
}

func (request KubernetesResourcesRequest) ServiceWithPort(serviceName string, port int) HttpGet {
	service := k8s.GetService(request.t, request.options, serviceName)
	url := fmt.Sprintf("http://%s", k8s.GetServiceEndpoint(request.t, request.options, service, port))
	return HttpGet{url:url, t:request.t}
}

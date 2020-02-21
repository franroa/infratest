package test_cases

import (
	"fmt"
	"github.com/franroa/infratest/environment"
	"github.com/franroa/infratest/helm"
	http_helper "github.com/franroa/infratest/http-helper"
	"github.com/franroa/infratest/k8s"
	"testing"
	"time"
)

func TestThatEksClusterIsAvailable(t *testing.T) {
	t.Parallel()
	environment.SetPropertyFile("../env.yaml")
	helmChart := "./../charts/hello-world"
	releaseName := "hello-world"
	kubectlOptions := k8s.NewKubectlOptions(t, "", "default")
	defer helm.RunHelmUninstall(t, releaseName)
	helm.InstallWithKubectlOptions(t, kubectlOptions, helmChart, releaseName)

	k8s.WaitUntilServiceAvailable(t, kubectlOptions, "hello-world", 10, 1*time.Second)
	service := k8s.GetService(t, kubectlOptions, "hello-world")
	url := fmt.Sprintf("http://%s", k8s.GetServiceEndpoint(t, kubectlOptions, service, 80))

	//http_helper.HttpGetWithRetry(t, url, nil, 200, "Hello world!", 50, 30*time.Second)
	http_helper.HttpGetWithRetry(t, url, nil, 200, "Hello world!", 50, 3*time.Second)

	// TODO -> If the response has been retrieved and the body or the status are not the expected ones, the retry should stop
}

//func TestThatEksClusterIsAvailableWorkingWithNamespaces(t *testing.T) {
//	t.Parallel()
//
//	helmChart := "./../examples/hello-world"
//	releaseName := "hello-world"
//	namespaceName := "hello-world"
//	kubectlOptions := k8ssynaos.NewKubectlOptions(t, "", namespaceName)
//
//	defer k8ssynaos.DeleteNamespace(t, kubectlOptions, namespaceName)
//	k8ssynaos.CreateNamespace(t, kubectlOptions, namespaceName)
//
//	//defer helm.RunHelmUninstallE(t, releaseName)
//	// Override service type to node port
//	options := &helm.Options{
//		KubectlOptions: kubectlOptions,
//		SetValues: map[string]string{
//			"service.type": "NodePort",
//		},
//	}
//	helm.Install(t, options, helmChart, releaseName)
//
//
//	// website::tag::4:: Verify the service is available and get the URL for it.
//	k8ssynaos.WaitUntilServiceAvailable(t, kubectlOptions, "hello-world", 10, 1*time.Second)
//	service := k8ssynaos.GetService(t, kubectlOptions, "hello-world")
//	url := fmt.Sprintf("http://%s", k8ssynaos.GetServiceEndpoint(t, kubectlOptions, service, 5000))
//
//	// website::tag::5:: Make an HTTP request to the URL and make sure it returns a 200 OK with the body "Hello, World".
//	http_helper.HttpGetWithRetry(t, url, nil, 200, "Hello world!", 30, 3*time.Second)
//}

package k8s

import (
	"github.com/stretchr/testify/require"
	"github.com/franroa/infratest/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

// GetNodes queries Kubernetes for information about the worker nodes registered to the cluster. If anything goes wrong,
// the function will automatically fail the test.
func GetNodes(t *testing.T, options *KubectlOptions) []corev1.Node {
	nodes, err := GetNodesE(t, options)
	require.NoError(t, err)
	return nodes
}

// GetNodesE queries Kubernetes for information about the worker nodes registered to the cluster.
func GetNodesE(t *testing.T, options *KubectlOptions) ([]corev1.Node, error) {
	return GetNodesByFilterE(t, options, metav1.ListOptions{})
}

// GetNodesByFilterE queries Kubernetes for information about the worker nodes registered to the cluster, filtering the
// list of nodes using the provided ListOptions.
func GetNodesByFilterE(t *testing.T, options *KubectlOptions, filter metav1.ListOptions) ([]corev1.Node, error) {
	logger.Logf(t, "Getting list of nodes from Kubernetes")

	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}

	nodes, err := clientset.CoreV1().Nodes().List(filter)
	if err != nil {
		return nil, err
	}
	return nodes.Items, err
}

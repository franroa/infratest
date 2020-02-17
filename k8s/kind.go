package k8s

import (
	"strings"
	"testing"
)

// IsMinikubeE returns true if the underlying kubernetes cluster is Minikube. This is determined by getting the
// associated nodes and checking if:
// - there is only one node
// - the node is named "minikube"
func IsKindWorkerE(t *testing.T, options *KubectlOptions) (bool, error) {
	nodes, err := GetNodesE(t, options)
	if err != nil {
		return false, err
	}
	return len(nodes) > 0 && (strings.HasPrefix(nodes[0].Name, "kind-worker") || strings.HasPrefix(nodes[0].Name, "kind-control-plane")), nil
}

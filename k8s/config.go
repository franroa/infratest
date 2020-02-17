package k8s

import (
	"github.com/franroa/infratest/environment"
	"github.com/mitchellh/go-homedir"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"path/filepath"
	"testing"
)

// LoadApiClientConfigE will load a ClientConfig object from a file path that points to a location on disk containing a
// kubectl config, with the requested context loaded.
func LoadApiClientConfigE(configPath string, contextName string) (*restclient.Config, error) {
	overrides := clientcmd.ConfigOverrides{}
	if contextName != "" {
		overrides.CurrentContext = contextName
	}
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: configPath},
		&overrides)
	return config.ClientConfig()
}

// GetKubeConfigPathE determines which file path to use as the kubectl config path
func GetKubeConfigPathE(t *testing.T) (string, error) {
	kubeConfigPath := environment.GetFirstNonEmptyEnvVarOrEmptyString(t, []string{"KUBECONFIG"})
	if kubeConfigPath == "" {
		configPath, err := KubeConfigPathFromHomeDirE()
		if err != nil {
			return "", err
		}
		kubeConfigPath = configPath
	}
	return kubeConfigPath, nil
}

// KubeConfigPathFromHomeDirE returns a string to the default Kubernetes config path in the home directory. This will
// error if the home directory can not be determined.
func KubeConfigPathFromHomeDirE() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(home, ".kube", "config")
	return configPath, err
}
package helm

import (
	"github.com/franroa/infratest/environment"
	"github.com/stretchr/testify/require"
	"github.com/franroa/infratest/shell"
	"testing"
)

func RunHelmUninstall(t *testing.T, releaseName string) {
	require.NoError(t, RunHelmUninstallE(t, releaseName))
}

func RunHelmUninstallE(t *testing.T, release string) error {
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, "--kube-context", environment.GetContextByHelmEnvironment(t), "uninstall", release)
	command := shell.Command{
		Command: "helm",
		Args:    cmdArgs,
	}
	_, err := shell.RunCommandAndGetOutputE(t, command)
	return err
}

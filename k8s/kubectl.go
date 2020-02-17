package k8s

import (
	"github.com/franroa/infratest/shell"
	"testing"

	"github.com/stretchr/testify/require"
)

// RunKubectl will call kubectl using the provided options and args, failing the test on error.
func RunKubectl(t *testing.T, options *KubectlOptions, args ...string) {
	require.NoError(t, RunKubectlE(t, options, args...))
}

// RunKubectlE will call kubectl using the provided options and args.
func RunKubectlE(t *testing.T, options *KubectlOptions, args ...string) error {
	_, err := RunKubectlAndGetOutputE(t, options, args...)
	return err
}

// RunKubectlAndGetOutputE will call kubectl using the provided options and args, returning the output of stdout and
// stderr.
func RunKubectlAndGetOutputE(t *testing.T, options *KubectlOptions, args ...string) (string, error) {
	cmdArgs := []string{}
	if options.ContextName != "" {
		cmdArgs = append(cmdArgs, "--context", options.ContextName)
	}
	if options.ConfigPath != "" {
		cmdArgs = append(cmdArgs, "--kubeconfig", options.ConfigPath)
	}
	if options.Namespace != "" {
		cmdArgs = append(cmdArgs, "--namespace", options.Namespace)
	}
	cmdArgs = append(cmdArgs, args...)
	command := shell.Command{
		Command: "kubectl",
		Args:    cmdArgs,
		Env:     options.Env,
	}
	return shell.RunCommandAndGetOutputE(t, command)
}

// KubectlDelete will take in a file path and delete it from the cluster targeted by KubectlOptions. If there are any
// errors, fail the test immediately.
func KubectlDelete(t *testing.T, options *KubectlOptions, configPath string) {
	require.NoError(t, KubectlDeleteE(t, options, configPath))
}

// KubectlDeleteE will take in a file path and delete it from the cluster targeted by KubectlOptions.
func KubectlDeleteE(t *testing.T, options *KubectlOptions, configPath string) error {
	return RunKubectlE(t, options, "delete", "-f", configPath)
}


// KubectlApply will take in a file path and apply it to the cluster targeted by KubectlOptions. If there are any
// errors, fail the test immediately.
func KubectlApply(t *testing.T, options *KubectlOptions, configPath string) {
	require.NoError(t, KubectlApplyE(t, options, configPath))
}

// KubectlApplyE will take in a file path and apply it to the cluster targeted by KubectlOptions.
func KubectlApplyE(t *testing.T, options *KubectlOptions, configPath string) error {
	return RunKubectlE(t, options, "apply", "-f", configPath)
}


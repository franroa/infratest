package environment

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var properties = Properties{}

type Environment struct {
	ContextName    string `yaml:"context_name"`
	ValuesFileName string `yaml:"values_file"`
}

type Properties struct {
	DefaultEnvironment string                 `yaml:"default_environment"`
	HelmEnvKey         string                 `helm_environment_key`
	Environments       map[string]Environment `yaml:"environments"`
}

func SetPropertyFile(propertiesFile string) {
	filename, _ := filepath.Abs(propertiesFile)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &properties)
	if err != nil {
		panic(err)
	}
}

func GetProperties() Properties {
	return properties
}

func GetValueFileByEnvironment(t *testing.T, chart string) string {
	properties := GetProperties()
	helmEnv := GetFirstNonEmptyEnvVarOrEmptyString(t, []string{properties.HelmEnvKey})

	if helmEnv == "" {
		helmEnv = properties.DefaultEnvironment
	}

	return chart + "/" + properties.Environments[helmEnv].ValuesFileName
}

func GetContextByHelmEnvironment(t *testing.T) string {
	properties := GetProperties()
	helmEnv := GetFirstNonEmptyEnvVarOrEmptyString(t, []string{properties.HelmEnvKey})

	if helmEnv == "" {
		helmEnv = properties.DefaultEnvironment
	}

	return properties.Environments[helmEnv].ContextName
}

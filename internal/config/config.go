package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
)

// DirectoryConfig represents the configuration for a single monitored directory
type DirectoryConfig struct {
	LogFile       string
	IncludeRegexp []*regexp.Regexp
	ExcludeRegexp []*regexp.Regexp
	Commands      []string
}

// Config represents the top-level configuration object
type Config struct {
	Directories map[string]DirectoryConfig
}

// Reader interface defines the behavior for reading the configuration from a file
type Reader interface {
	ReadFile(path string) ([]byte, error)
}

// DefaultReader implements the Reader interface using os.ReadFile and yaml.Unmarshal
type DefaultReader struct{}

func (dcr *DefaultReader) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func ParseConfig(content []byte) (*Config, error) {
	if len(content) == 0 {
		return nil, fmt.Errorf("empty config file")
	}
	var configs []struct {
		Path          string   `yaml:"path"`
		LogFile       string   `yaml:"log_file"`
		IncludeRegexp []string `yaml:"include_regexp"`
		ExcludeRegexp []string `yaml:"exclude_regexp"`
		Commands      []string `yaml:"commands"`
	}
	mapConfigs := map[string]DirectoryConfig{}
	err := yaml.Unmarshal(content, &configs)
	if err != nil {
		return nil, err
	}

	for _, v := range configs {
		if v.Path == "" {
			return nil, fmt.Errorf("path in config file can not be empty")
		}
		if len(v.Commands) == 0 {
			return nil, fmt.Errorf("no commands for path: %s", v.Path)
		}

		IncludeRegexp, err := getRegexpArr(v.IncludeRegexp)
		if err != nil {

		}

		ExcludeRegexp, err := getRegexpArr(v.ExcludeRegexp)
		if err != nil {

		}

		mapConfigs[v.Path] = DirectoryConfig{
			LogFile:       v.LogFile,
			IncludeRegexp: IncludeRegexp,
			ExcludeRegexp: ExcludeRegexp,
			Commands:      v.Commands,
		}
	}

	return &Config{Directories: mapConfigs}, err
}

// ReadConfig reads the configuration from a file at the given path using the provided ConfigReader
func ReadConfig(path string, reader Reader) (*Config, error) {
	content, err := reader.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %s", err)
	}

	config, err := ParseConfig(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file, %s", err)
	}

	return config, nil
}

func getRegexpArr(expr []string) ([]*regexp.Regexp, error) {
	ret := make([]*regexp.Regexp, len(expr))
	for i := range expr {
		reg, err := regexp.Compile(expr[i])
		if err != nil {
			return nil, err
		}
		ret[i] = reg
	}
	return ret, nil
}

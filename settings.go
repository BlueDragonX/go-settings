package settings

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

type Settings struct {
	Key    string
	Values map[interface{}]interface{}
}

// Parse the provided YAML into a new Settings object.
func Parse(data []byte) (*Settings, error) {
	values := make(map[interface{}]interface{})
	if err := yaml.Unmarshal(data, values); err == nil {
		return &Settings{Values: values}, nil
	} else {
		return nil, err
	}
}

// Read and parse settings from the provided reader.
func Read(reader io.Reader) (*Settings, error) {
	if data, err := ioutil.ReadAll(reader); err == nil {
		return Parse(data)
	} else {
		return nil, err
	}
}

// Load and parse settings from the file at the provided path.
func Load(path string) (*Settings, error) {
	if file, err := os.Open(path); err == nil {
		defer file.Close()
		return Read(file)
	} else {
		return nil, err
	}
}

// Load and parse settings from the file at the provided path. If an error
// occurs print it to stderr and call os.Exit(1).
func LoadOrExit(path string) *Settings {
	settings, err := Load(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return settings
}

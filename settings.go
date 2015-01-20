package settings

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var KeyError error = errors.New("key not found")
var TypeError error = errors.New("invalid type conversion")

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

// Get a value from the settings object.
func (s *Settings) Raw(key string) (interface{}, error) {
	names := strings.Split(key, ".")
	var data interface{} = s.Values
	for _, name := range names {
		if items, ok := data.(map[interface{}]interface{}); ok {
			if data, ok = items[name]; !ok {
				return nil, KeyError
			}
		} else if items, ok := data.([]interface{}); ok {
			if n, err := strconv.Atoi(name); err == nil {
				if n < len(items) {
					data = items[n]
				} else {
					return nil, KeyError
				}
			} else {
				return nil, KeyError
			}
		} else {
			return nil, TypeError
		}
	}
	return data, nil
}

// Get a settings object.
func (s *Settings) Object(key string) (*Settings, error) {
	if value, err := s.Raw(key); err == nil {
		if mapping, ok := value.(map[interface{}]interface{}); ok {
			return &Settings{Key: key, Values: mapping}, nil
		} else {
			return nil, TypeError
		}
	} else {
		return nil, err
	}
}

// Get an array of settings objects.
func (s *Settings) ObjectArray(key string) ([]*Settings, error) {
	if value, err := s.Raw(key); err == nil {
		if items, ok := value.([]interface{}); ok {
			array := make([]*Settings, len(items))
			for n, item := range items {
				if mapping, ok := item.(map[interface{}]interface{}); ok {
					settingsKey := fmt.Sprintf("%s.%d", key, n)
					array[n] = &Settings{Key: settingsKey, Values: mapping}
				} else {
					return nil, TypeError
				}
			}
			return array, nil
		} else {
			return nil, TypeError
		}
	} else {
		return nil, err
	}
}

// Get a string value.
func (s *Settings) String(key string) (string, error) {
	if value, err := s.Raw(key); err == nil {
		if valueStr, ok := value.(string); ok {
			return valueStr, nil
		} else {
			return "", TypeError
		}
	} else {
		return "", err
	}
}

// Get an array of string values.
func (s *Settings) StringArray(key string) ([]string, error) {
	if value, err := s.Raw(key); err == nil {
		if items, ok := value.([]interface{}); ok {
			array := make([]string, len(items))
			for n, item := range items {
				switch item.(type) {
				case string:
					array[n] = item.(string)
				default:
					return nil, TypeError
				}
			}
			return array, nil
		} else {
			return nil, TypeError
		}
	} else {
		return nil, err
	}
}

// Get an integer value.
func (s *Settings) Int(key string) (int, error) {
	if value, err := s.Raw(key); err == nil {
		if valueInt, ok := value.(int); ok {
			return valueInt, nil
		} else {
			return 0, TypeError
		}
	} else {
		return 0, err
	}
}

// Get an array of integer values.
func (s *Settings) IntArray(key string) ([]int, error) {
	if value, err := s.Raw(key); err == nil {
		if items, ok := value.([]interface{}); ok {
			array := make([]int, len(items))
			for n, item := range items {
				switch item.(type) {
				case int:
					array[n] = item.(int)
				default:
					return nil, TypeError
				}
			}
			return array, nil
		} else {
			return nil, TypeError
		}
	} else {
		return nil, err
	}
}

// Get a float value.
func (s *Settings) Float(key string) (float64, error) {
	if value, err := s.Raw(key); err == nil {
		switch value.(type) {
		case float64:
			return value.(float64), nil
		case int:
			return float64(value.(int)), nil
		default:
			return 0, TypeError
		}
	} else {
		return 0, err
	}
}

// Get an array of float values.
func (s *Settings) FloatArray(key string) ([]float64, error) {
	if value, err := s.Raw(key); err == nil {
		if items, ok := value.([]interface{}); ok {
			array := make([]float64, len(items))
			for n, item := range items {
				switch item.(type) {
				case float64:
					array[n] = item.(float64)
				case int:
					array[n] = float64(item.(int))
				default:
					return nil, TypeError
				}
			}
			return array, nil
		} else {
			return nil, TypeError
		}
	} else {
		return nil, err
	}
}

// Get a boolean value.
func (s *Settings) Bool(key string) (bool, error) {
	if value, err := s.Raw(key); err == nil {
		switch value.(type) {
		case bool:
			return value.(bool), nil
		case int:
			return value.(int) != 0, nil
		case float64:
			return value.(float64) != 0, nil
		case string:
			if valueBool, err := strconv.ParseBool(value.(string)); err == nil {
				return valueBool, nil
			} else {
				return false, TypeError
			}
		default:
			return false, TypeError
		}
	} else {
		return false, err
	}
}

// Get an array of boolean values.
func (s *Settings) BoolArray(key string) ([]bool, error) {
	if value, err := s.Raw(key); err == nil {
		if items, ok := value.([]interface{}); ok {
			array := make([]bool, len(items))
			for n, item := range items {
				switch item.(type) {
				case bool:
					array[n] = item.(bool)
				case int:
					array[n] = item.(int) != 0
				case float64:
					array[n] = item.(float64) != 0
				case string:
					if valueBool, err := strconv.ParseBool(item.(string)); err == nil {
						array[n] = valueBool
					} else {
						return nil, TypeError
					}
				default:
					return nil, TypeError
				}
			}
			return array, nil
		} else {
			return nil, TypeError
		}
	} else {
		return nil, err
	}
}

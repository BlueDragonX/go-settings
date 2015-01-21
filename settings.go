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

var IndexError error = errors.New("invalid index")
var KeyError error = errors.New("key not found")
var ObjectError error = errors.New("invalid object")
var RangeError error = errors.New("index out of range")
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

// Get the value at the index of the provided map or array.
func getElement(obj interface{}, index string) (interface{}, error) {
	switch obj.(type) {
	case map[interface{}]interface{}:
		if item, ok := (obj.(map[interface{}]interface{}))[index]; ok {
			return item, nil
		} else {
			return nil, IndexError
		}
	case []interface{}:
		if n, err := strconv.Atoi(index); err == nil {
			return (obj.([]interface{}))[n], nil
		}
	}
	return nil, ObjectError
}

// Set the value at the index of the provided map or array.
func setElement(obj interface{}, index string, value interface{}) error {
	if obj, ok := value.(*Settings); ok {
		value = obj.Values
	}

	if mapping, ok := obj.(map[interface{}]interface{}); ok {
		mapping[index] = value
	} else if array, ok := obj.([]interface{}); ok {
		n, err := strconv.Atoi(index)
		if err != nil {
			return IndexError
		}
		if n < 0 || n >= len(array) {
			return RangeError
		}
		array[n] = value
	} else {
		return ObjectError
	}
	return nil
}

// Create maps in `values` along the provided path and return the last created map.
func createPath(values interface{}, path []string) (interface{}, error) {
	for _, name := range path {
		if next, err := getElement(values, name); err == nil {
			values = next
		} else if err == IndexError {
			next = make(map[interface{}]interface{})
			setElement(values, name, next)
			values = next
		} else {
			return nil, err
		}
	}
	return values, nil
}

// Get the parent object for the provided path creating any missing elements as necessary.
func getParent(value interface{}, path []string) (parent interface{}, err error) {
	if len(path) == 1 {
		parent = value
	} else {
		parent, err = createPath(value, path[:len(path)-1])
	}
	return
}

// Set a value in the settings object. This overrides all keys in the path that
// are not already objects. The following errors may be returned:
//
// ObjectError - a child object is not a `map[interface{}]interface{}` or `[]interface{}`
// IndexError - a key cannot be converted to an integer for a child array
// RangeError - the index is out of range for a child array
func (s *Settings) Set(key string, value interface{}) error {
	names := strings.Split(key, ".")
	if parent, err := createPath(s.Values, names[:len(names)-1]); err == nil {
		return setElement(parent, names[len(names)-1], value)
	} else {
		return err
	}
}

// Append a value to an array. Creates an array at that location if it does not
// exist. The following errors may be returned:
//
// ObjectError - a child object is not a `map[interface{}]interface{}` or `[]interface{}`
// IndexError - a key cannot be converted to an integer for a child array
func (s *Settings) Append(key string, value interface{}) error {
	var ok bool
	var err error
	var parent, arrayObj interface{}
	var array []interface{}

	if obj, ok := value.(*Settings); ok {
		value = obj.Values
	}

	names := strings.Split(key, ".")
	name := names[len(names)-1]
	if parent, err = getParent(s.Values, names); err != nil {
		return err
	}

	if arrayObj, err = getElement(parent, name); err == IndexError {
		array = make([]interface{}, 0)
	} else if err != nil {
		return err
	} else if array, ok = arrayObj.([]interface{}); !ok {
		array = make([]interface{}, 0)
	}

	return setElement(parent, name, append(array, value))
}

// Delete a key. May return an error on failure. A non-existent key is not an error case.
//
// ObjectError - a child object is not a `map[interface{}]interface{}` or `[]interface{}`
// IndexError - a key cannot be converted to an integer for a child array
func (s *Settings) Delete(key string) error {
	var err error

	names := strings.Split(key, ".")
	name := names[len(names)-1]
	if len(names) == 1 {
		if _, ok := s.Values[name]; ok {
			delete(s.Values, name)
		}
	} else {
		var child, parent interface{}
		path := names[:len(names)-1]
		childName := names[len(path)-1]
		if parent, err = getParent(s.Values, path); err != nil {
			return err
		}
		if child, err = getElement(parent, childName); err != nil {
			return err
		}

		switch child.(type) {
		case map[interface{}]interface{}:
			delete(child.(map[interface{}]interface{}), name)
		case []interface{}:
			if n, err := strconv.Atoi(name); err == nil {
				array := child.([]interface{})
				err = setElement(parent, childName, append(array[:n], array[n+1:]...))
			} else {
				err = IndexError
			}
		default:
			err = ObjectError
		}
	}
	return err
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

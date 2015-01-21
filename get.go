package settings

import (
	"fmt"
	"strconv"
	"strings"
)

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

// Get a map of settings objects.
func (s *Settings) ObjectMap(key string) (map[string]*Settings, error) {
	if raw, err := s.Raw(key); err == nil {
		rawMap, ok := raw.(map[interface{}]interface{})
		if !ok {
			return nil, TypeError
		}

		objectMap := make(map[string]*Settings)
		for rawMapKey, rawMapValue := range rawMap {
			if settingsValues, ok := rawMapValue.(map[interface{}]interface{}); ok {
				keyStr := fmt.Sprintf("%v", rawMapKey)
				settingsKey := key + "." + keyStr
				objectMap[keyStr] = &Settings{Key: settingsKey, Values: settingsValues}
			} else {
				return nil, TypeError
			}
		}
		return objectMap, nil
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

// Get a map of strings.
func (s *Settings) StringMap(key string) (map[string]string, error) {
	if raw, err := s.Raw(key); err == nil {
		rawMap, ok := raw.(map[interface{}]interface{})
		if !ok {
			return nil, TypeError
		}

		stringMap := make(map[string]string)
		for rawMapKey, rawMapValue := range rawMap {
			if stringValue, ok := rawMapValue.(string); ok {
				keyStr := fmt.Sprintf("%v", rawMapKey)
				stringMap[keyStr] = stringValue
			} else {
				return nil, TypeError
			}
		}
		return stringMap, nil
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

// Get a map of integers.
func (s *Settings) IntMap(key string) (map[string]int, error) {
	if raw, err := s.Raw(key); err == nil {
		rawMap, ok := raw.(map[interface{}]interface{})
		if !ok {
			return nil, TypeError
		}

		intMap := make(map[string]int)
		for rawMapKey, rawMapValue := range rawMap {
			if intValue, ok := rawMapValue.(int); ok {
				keyStr := fmt.Sprintf("%v", rawMapKey)
				intMap[keyStr] = intValue
			} else {
				return nil, TypeError
			}
		}
		return intMap, nil
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

// Get a map of floats.
func (s *Settings) FloatMap(key string) (map[string]float64, error) {
	if raw, err := s.Raw(key); err == nil {
		rawMap, ok := raw.(map[interface{}]interface{})
		if !ok {
			return nil, TypeError
		}

		floatMap := make(map[string]float64)
		for rawMapKey, rawMapValue := range rawMap {
			if intValue, ok := rawMapValue.(int); ok {
				keyStr := fmt.Sprintf("%v", rawMapKey)
				floatMap[keyStr] = float64(intValue)
			} else if floatValue, ok := rawMapValue.(float64); ok {
				keyStr := fmt.Sprintf("%v", rawMapKey)
				floatMap[keyStr] = floatValue
			} else {
				return nil, TypeError
			}
		}
		return floatMap, nil
	} else {
		return nil, err
	}
}

// Convert a value to a bool.
func getBoolValue(value interface{}) (bool, error) {
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
}

// Get a boolean value.
func (s *Settings) Bool(key string) (bool, error) {
	if value, err := s.Raw(key); err == nil {
		return getBoolValue(value)
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
				if boolValue, err := getBoolValue(item); err == nil {
					array[n] = boolValue
				} else {
					return nil, err
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

// Get a map of bools.
func (s *Settings) BoolMap(key string) (map[string]bool, error) {
	if raw, err := s.Raw(key); err == nil {
		rawMap, ok := raw.(map[interface{}]interface{})
		if !ok {
			return nil, TypeError
		}

		boolMap := make(map[string]bool)
		for rawMapKey, rawMapValue := range rawMap {
			keyStr := fmt.Sprintf("%v", rawMapKey)
			if boolValue, err := getBoolValue(rawMapValue); err == nil {
				boolMap[keyStr] = boolValue
			} else {
				return nil, err
			}
		}
		return boolMap, nil
	} else {
		return nil, err
	}
}

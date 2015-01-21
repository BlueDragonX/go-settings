package settings

import (
	"strconv"
	"strings"
)

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
	if s.Values == nil {
		s.Values = make(map[interface{}]interface{})
	}
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

	if s.Values == nil {
		s.Values = make(map[interface{}]interface{})
	}

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

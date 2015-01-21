package settings

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)


func isElementEqual(t *testing.T, a, b interface{}) bool {
	if reflect.TypeOf(a) == reflect.TypeOf(b) {
		if arrayA, ok := a.([]interface{}); ok {
			arrayB := b.([]interface{})
			return isArrayEqual(t, arrayA, arrayB)
		} else if mapA, ok := a.(map[interface{}]interface{}); ok {
			mapB := b.(map[interface{}]interface{})
			return isMapEqual(t, mapA, mapB)
		} else {
			if a != b {
				t.Log("values unequal")
			}
			return a == b
		}
	}
	t.Logf("unmatched types %v, %v\n", reflect.TypeOf(a), reflect.TypeOf(b))
	return false
}

func isArrayEqual(t *testing.T, a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for n := range a {
		if !isElementEqual(t, a[n], b[n]) {
			t.Logf("index %d unequal\n", n)
			return false
		}
	}
	return true
}

func isMapEqual(t *testing.T, a, b map[interface{}]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if !isElementEqual(t, v, b[k]) {
			t.Logf("key %v unequal\n", k)
			return false
		}
	}
	return true
}

func getSettings() *Settings {
	settings, _ := Parse([]byte(`key: value

mapping:
  a: aye
  b: bee

string-array:
- one
- two

integer-array:
- 1
- 2

float-array:
- 1.3
- 2.2
- 3.1

bool-array:
- true
- true
- false
- true

settings-array:
- name: one
  value: I won!
- name: two
  value: Me too!

mixed-array:
- one
- 2

values:
  bool: true
  integer: 1
  float: 2.3
  string: value

nested:
  array:
  - one
  - two`))
	return settings
}

func getBasicInput() ([]byte, map[interface{}]interface{}) {
	data := `a: aye
b:
  c: see
  d: dee
e:
- one
- two`

	want := make(map[interface{}]interface{})
	want["a"] = "aye"
	mapping := make(map[interface{}]interface{})
	mapping["c"] = "see"
	mapping["d"] = "dee"
	want["b"] = mapping
	array := make([]interface{}, 2)
	array[0] = "one"
	array[1] = "two"
	want["e"] = array

	return []byte(data), want
}

func getBasicFile(t *testing.T) (string, map[interface{}]interface{}) {
	data, want := getBasicInput()
	file, err := ioutil.TempFile("", "go-settings-")
	if err != nil {
		t.Fatal("failed to create temp file for testing")
	}
	defer func() {
		file.Close()
	}()

	path := file.Name()
	file.Write(data)
	file.Close()
	return path, want
}

func TestParse(t *testing.T) {
	data, want := getBasicInput()
	if have, err := Parse(data); err == nil {
		if !isElementEqual(t, have.Values, want) {
			t.Errorf("%v != %v", want, have)
		}
	} else {
		t.Error(err)
	}
}

func TestRead(t *testing.T) {
	data, want := getBasicInput()
	reader := bytes.NewBuffer(data)
	if have, err := Read(reader); err == nil {
		if !isElementEqual(t, have.Values, want) {
			t.Errorf("%v != %v", want, have)
		}
	} else {
		t.Error(err)
	}
}

func TestLoad(t *testing.T) {
	path, want := getBasicFile(t)
	defer os.Remove(path)

	if have, err := Load(path); err == nil {
		if !isElementEqual(t, have.Values, want) {
			t.Errorf("%v != %v", want, have)
		}
	} else {
		t.Error(err)
	}
}

func TestLoadOrExit(t *testing.T) {
	path, want := getBasicFile(t)
	defer os.Remove(path)

	have := LoadOrExit(path)
	if !isElementEqual(t, have.Values, want) {
		t.Errorf("%v != %v", want, have)
	}
}

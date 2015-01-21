package settings

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func isSettingsArrayEqual(t *testing.T, a, b []*Settings) bool {
	if len(a) != len(b) {
		return false
	}
	for n := range a {
		if !reflect.DeepEqual(a[n].Values, b[n].Values) {
			t.Logf("index %d unequal: %v != %v\n", n, a, b)
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

string-map:
  a: aye
  b: bee

integer-array:
- 1
- 2

integer-map:
  one: 1
  two: 2

float-array:
- 1.3
- 2.2
- 3.1

float-map:
  one: 1.1
  two: 2.2

bool-array:
- true
- true
- false
- true

bool-map:
  "yes": true
  "no": false

settings-array:
- name: one
  value: I won!
- name: two
  value: Me too!

settings-map:
  one:
    value: I won!
  two:
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
		if !reflect.DeepEqual(have.Values, want) {
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
		if !reflect.DeepEqual(have.Values, want) {
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
		if !reflect.DeepEqual(have.Values, want) {
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
	if !reflect.DeepEqual(have.Values, want) {
		t.Errorf("%v != %v", want, have)
	}
}

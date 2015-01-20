package settings

import (
	"testing"
)

var input string = `key: value

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
  string: value`

func checkMapEqual(t *testing.T, want, have map[interface{}]interface{}) bool {
	if len(want) != len(have) {
		t.Errorf("%v != %v", want, have)
		return false
	}
	for k, v := range want {
		if v != have[k] {
			t.Errorf("want[%v] != have[%v] (%v != %v)", k, k, v, have[k])
			return false
		}
	}
	return true
}

func checkIsMap(t *testing.T, have interface{}) bool {
	if _, ok := have.(map[interface{}]interface{}); !ok {
		t.Errorf("%v is not a map", have)
		return false
	}
	return true
}

func checkIsArray(t *testing.T, have interface{}) bool {
	if _, ok := have.([]interface{}); !ok {
		t.Errorf("%v is not an array", have)
		return false
	}
	return true
}

func checkSettings(t *testing.T, want *Settings, have interface{}) bool {
	if settings, ok := have.(*Settings); ok {
		if settings.Key != want.Key {
			t.Errorf("key is incorrect (%s != %s)", want.Key, settings.Key)
		}
		checkMapEqual(t, settings.Values, want.Values)
	} else {
		t.Errorf("%v is not a *Settings", have)
		return false
	}
	return true
}

func checkBool(t *testing.T, want bool, have interface{}) bool {
	if boolValue, ok := have.(bool); ok {
		if boolValue != want {
			t.Errorf("%t != %t", want, boolValue)
			return false
		}
	} else {
		t.Errorf("%v is not a bool", have)
		return false
	}
	return true
}

func checkInt(t *testing.T, want int, have interface{}) bool {
	if intValue, ok := have.(int); ok {
		if intValue != want {
			t.Errorf("%d != %d", want, intValue)
			return false
		}
	} else {
		t.Errorf("%v is not an integer", have)
		return false
	}
	return true
}

func checkFloat(t *testing.T, want float64, have interface{}) bool {
	if floatValue, ok := have.(float64); ok {
		if floatValue != want {
			t.Errorf("%f != %f", want, floatValue)
			return false
		}
	} else {
		t.Errorf("%v is not an float", have)
		return false
	}
	return true
}

func checkString(t *testing.T, want string, have interface{}) bool {
	if stringValue, ok := have.(string); ok {
		if stringValue != want {
			t.Errorf("%s != %s", want, stringValue)
			return false
		}
	} else {
		t.Errorf("%v is not a string", have)
		return false
	}
	return true
}

func TestRaw(t *testing.T) {
	var key string
	settings, _ := Parse([]byte(input))
	t.Logf("%v\n", settings.Values)

	// retrieve a map
	if value, err := settings.Raw("mapping"); err == nil {
		checkIsMap(t, value)
	} else {
		t.Error(err)
	}

	// retrieve an array
	if value, err := settings.Raw("string-array"); err == nil {
		checkIsArray(t, value)
	} else {
		t.Error(err)
	}

	// retrieve a single value
	if value, err := settings.Raw("key"); err == nil {
		checkString(t, "value", value)
	} else {
		t.Error(err)
	}

	// retrieve a child of a map
	if value, err := settings.Raw("mapping.a"); err == nil {
		checkString(t, "aye", value)
	} else {
		t.Error(err)
	}

	// retrieve a child of an array
	if value, err := settings.Raw("string-array.1"); err == nil {
		checkString(t, "two", value)
	} else {
		t.Error(err)
	}

	// retrieve missing value
	key = "sir-not-appearing-in-this-film"
	if _, err := settings.Raw(key); err != KeyError {
		t.Errorf("%s did not cause a KeyError", key)
	}

	// retrieve missing map value
	key = "mapping.c"
	if _, err := settings.Raw(key); err != KeyError {
		t.Errorf("%s did not cause a KeyError", key)
	}

	// retrieve missing array value
	key = "string-array.3"
	if _, err := settings.Raw(key); err != KeyError {
		t.Errorf("%s did not cause a KeyError", key)
	}

	// retrieve a bool
	key = "values.bool"
	if value, err := settings.Raw(key); err == nil {
		checkBool(t, true, value)
	} else {
		t.Error(err)
	}

	// retrieve an int
	key = "values.integer"
	if value, err := settings.Raw(key); err == nil {
		checkInt(t, 1, value)
	} else {
		t.Error(err)
	}

	// retrieve a float
	key = "values.float"
	if value, err := settings.Raw(key); err == nil {
		checkFloat(t, 2.3, value)
	} else {
		t.Error(err)
	}

	// retrieve a string
	key = "values.string"
	if value, err := settings.Raw(key); err == nil {
		checkString(t, "value", value)
	} else {
		t.Error(err)
	}
}

func TestObject(t *testing.T) {
	var key string
	settings, _ := Parse([]byte(input))

	// test settings retrieval
	key = "values"
	if item, err := settings.Object(key); err == nil {
		values := make(map[interface{}]interface{})
		values["bool"] = true
		values["integer"] = 1
		values["float"] = 2.3
		values["string"] = "value"
		want := &Settings{Key: key, Values: values}
		checkSettings(t, want, item)
	} else {
		t.Error(err)
	}

	// test missing value
	key = "missing"
	if _, err := settings.Object(key); err != KeyError {
		t.Errorf("key %s found", key)
	}

	// test invalid value
	key = "string-array"
	if _, err := settings.Object(key); err != TypeError {
		t.Errorf("key %s is valid", key)
	}
}

func TestObjectArray(t *testing.T) {
	var key string
	settings, _ := Parse([]byte(input))

	// check valid
	key = "settings-array"
	if items, err := settings.ObjectArray(key); err == nil {
		want := make([]*Settings, 2)
		want1 := make(map[interface{}]interface{}, 2)
		want1["name"] = "one"
		want1["value"] = "I won!"
		want2 := make(map[interface{}]interface{}, 2)
		want2["name"] = "two"
		want2["value"] = "Me too!"
		want[0] = &Settings{Key: "settings-array.0", Values: want1}
		want[1] = &Settings{Key: "settings-array.1", Values: want2}
		if len(want) != len(items) {
			t.Errorf("settings array has incorrect length")
		}
		for n, item := range items {
			checkSettings(t, want[n], item)
		}
	} else {
		t.Error(err)
	}

	// check missing settings array
	key = "missing"
	if _, err := settings.ObjectArray(key); err != KeyError {
		t.Errorf("key %s found", key)
	}

	// check invalid type
	key = "string-array"
	if _, err := settings.ObjectArray(key); err != TypeError {
		t.Errorf("key %s is valid", key)
	}
}

func TestString(t *testing.T) {
	settings, _ := Parse([]byte(input))
	want := "value"
	if value, err := settings.String("values.string"); err == nil {
		if want != value {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}
}

func TestStringArray(t *testing.T) {
	settings, _ := Parse([]byte(input))

	// valid string array
	want := []string{"one", "two"}
	if value, err := settings.StringArray("string-array"); err == nil {
		if len(want) != len(value) {
			t.Errorf("%v != %v", want, value)
		}
		for n, item := range want {
			if item != value[n] {
				t.Errorf("%v != %v", want, value)
			}
		}
	} else {
		t.Error(err)
	}

	// mixed array
	key := "mixed-array"
	if _, err := settings.StringArray(key); err != TypeError {
		t.Errorf("key %s is valid", key)
	}
}

func TestInt(t *testing.T) {
	settings, _ := Parse([]byte(input))
	want := 1
	if value, err := settings.Int("values.integer"); err == nil {
		if want != value {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}
}

func TestIntArray(t *testing.T) {
	settings, _ := Parse([]byte(input))

	// valid string array
	want := []int{1, 2}
	if value, err := settings.IntArray("integer-array"); err == nil {
		if len(want) != len(value) {
			t.Errorf("%v != %v", want, value)
		}
		for n, item := range want {
			if item != value[n] {
				t.Errorf("%v != %v", want, value)
			}
		}
	} else {
		t.Error(err)
	}

	// mixed array
	key := "mixed-array"
	if _, err := settings.IntArray(key); err != TypeError {
		t.Errorf("key %s is valid", key)
	}
}

func TestFloat(t *testing.T) {
	settings, _ := Parse([]byte(input))
	want := 2.3
	if value, err := settings.Float("values.float"); err == nil {
		if want != value {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}
}

func TestFloatArray(t *testing.T) {
	settings, _ := Parse([]byte(input))

	// valid string array
	want := []float64{1.3, 2.2, 3.1}
	if value, err := settings.FloatArray("float-array"); err == nil {
		if len(want) != len(value) {
			t.Errorf("%v != %v", want, value)
		}
		for n, item := range want {
			if item != value[n] {
				t.Errorf("%v != %v", want, value)
			}
		}
	} else {
		t.Error(err)
	}

	// mixed array
	key := "mixed-array"
	if _, err := settings.FloatArray(key); err != TypeError {
		t.Errorf("key %s is valid", key)
	}
}

func TestBool(t *testing.T) {
	settings, _ := Parse([]byte(input))
	want := true
	if value, err := settings.Bool("values.bool"); err == nil {
		if want != value {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}
}

func TestBoolArray(t *testing.T) {
	settings, _ := Parse([]byte(input))

	// valid string array
	want := []bool{true, true, false, true}
	if value, err := settings.BoolArray("bool-array"); err == nil {
		if len(want) != len(value) {
			t.Errorf("%v != %v", want, value)
		}
		for n, item := range want {
			if item != value[n] {
				t.Errorf("%v != %v", want, value)
			}
		}
	} else {
		t.Error(err)
	}

	// mixed array
	key := "mixed-array"
	if _, err := settings.BoolArray(key); err != TypeError {
		t.Errorf("key %s is valid", key)
	}
}

package settings

import (
	"reflect"
	"testing"
	"time"
)

func TestRaw(t *testing.T) {
	// get a value from a new object
	settings := &Settings{}
	if _, err := settings.Raw("nope"); err != KeyError {
		t.Errorf("empty settings error is invalid: %s", err)
	}

	var key string
	settings = getSettings()

	// retrieve a map
	if value, err := settings.Raw("mapping"); err == nil {
		if _, ok := value.(map[interface{}]interface{}); !ok {
			t.Errorf("%v is not a map", value)
		}
	} else {
		t.Error(err)
	}

	// retrieve an array
	if value, err := settings.Raw("string-array"); err == nil {
		if _, ok := value.([]interface{}); !ok {
			t.Errorf("%v is not an array", value)
		}
	} else {
		t.Error(err)
	}

	// retrieve a single value
	if value, err := settings.Raw("key"); err == nil {
		if "value" != value {
			t.Errorf("%v != %v", "value", value)
		}
	} else {
		t.Error(err)
	}

	// retrieve a child of a map
	if value, err := settings.Raw("mapping.a"); err == nil {
		if "aye" != value {
			t.Errorf("%v != %v", "aye", value)
		}
	} else {
		t.Error(err)
	}

	// retrieve a child of an array
	if value, err := settings.Raw("string-array.1"); err == nil {
		if "two" != value {
			t.Errorf("%v != %v", "two", value)
		}
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
		if value != true {
			t.Errorf("%v != %v", true, value)
		}
	} else {
		t.Error(err)
	}

	// retrieve an int
	key = "values.integer"
	if value, err := settings.Raw(key); err == nil {
		if 1 != value {
			t.Errorf("%v != %v", 1, value)
		}
	} else {
		t.Error(err)
	}

	// retrieve a float
	key = "values.float"
	if value, err := settings.Raw(key); err == nil {
		if 2.3 != value {
			t.Errorf("%v != %v", 2.3, value)
		}
	} else {
		t.Error(err)
	}

	// retrieve a string
	key = "values.string"
	if value, err := settings.Raw(key); err == nil {
		if "value" != value {
			t.Errorf("%v != %v", "value", value)
		}
	} else {
		t.Error(err)
	}
}

func TestHas(t *testing.T) {
	type testInput struct {
		key string
		has bool
	}

	testInputs := []testInput{
		{"mapping", true},
		{"mapping.a", true},
		{"mapping.c", false},
		{"values", true},
		{"values.bool", true},
		{"values.bool.true", false},
		{"nope", false},
		{"nope.never", false},
	}

	settings := getSettings()
	for _, input := range testInputs {
		if settings.Has(input.key) != input.has {
			t.Errorf("Has(%s) != %t", input.key, input.has)
		}
	}
}

func TestObject(t *testing.T) {
	var key string
	settings := getSettings()

	// test settings retrieval
	key = "values"
	if item, err := settings.Object(key); err == nil {
		values := make(map[interface{}]interface{})
		values["bool"] = true
		values["integer"] = 1
		values["float"] = 2.3
		values["string"] = "value"
		values["duration"] = "5m"
		values["size"] = "15t"
		want := &Settings{Key: key, Values: values}
		if !reflect.DeepEqual(want, item) {
			t.Errorf("%v != %v", want, item)
		}
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
	settings := getSettings()

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

		if !reflect.DeepEqual(want, items) {
			t.Errorf("%v != %v", want, items)
		}
		for n, item := range items {
			if !reflect.DeepEqual(want[n], item) {
				t.Errorf("%v != %v", want[n], item)
			}
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

func TestObjectMap(t *testing.T) {
	settings := getSettings()

	want := make(map[string]*Settings)
	values := make(map[interface{}]interface{})
	values["value"] = "I won!"
	want["one"] = &Settings{Key: "settings-map.one", Values: values}
	values = make(map[interface{}]interface{})
	values["value"] = "Me too!"
	want["two"] = &Settings{Key: "settings-map.two", Values: values}

	if value, err := settings.ObjectMap("settings-map"); err == nil {
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}
}

func TestString(t *testing.T) {
	settings := getSettings()
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
	settings := getSettings()

	// valid string array
	want := []string{"one", "two"}
	if value, err := settings.StringArray("string-array"); err == nil {
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
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

func TestStringMap(t *testing.T) {
	settings := getSettings()

	// valid string map
	want := map[string]string{"a": "aye", "b": "bee"}
	if value, err := settings.StringMap("string-map"); err == nil {
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}
}

func TestInt(t *testing.T) {
	settings := getSettings()
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
	settings := getSettings()

	// valid string array
	want := []int{1, 2}
	if value, err := settings.IntArray("integer-array"); err == nil {
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
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

func TestIntMap(t *testing.T) {
	settings := getSettings()

	// valid string map
	want := map[string]int{"one": 1, "two": 2}
	if value, err := settings.IntMap("integer-map"); err == nil {
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}
}

func TestFloat(t *testing.T) {
	settings := getSettings()
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
	settings := getSettings()

	// valid string array
	want := []float64{1.3, 2.2, 3.1}
	if value, err := settings.FloatArray("float-array"); err == nil {
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
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

func TestFloatMap(t *testing.T) {
	settings := getSettings()

	// valid string map
	want := map[string]float64{"one": 1.1, "two": 2.2}
	if value, err := settings.FloatMap("float-map"); err == nil {
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}
}

func TestBool(t *testing.T) {
	settings := getSettings()
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
	settings := getSettings()

	// valid string array
	want := []bool{true, true, false, true}
	if value, err := settings.BoolArray("bool-array"); err == nil {
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
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

func TestBoolMap(t *testing.T) {
	settings := getSettings()

	// valid string map
	want := map[string]bool{"yes": true, "no": false}
	if value, err := settings.BoolMap("bool-map"); err == nil {
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}
}

func TestDuration(t *testing.T) {
	settings := getSettings()
	want := 5 * time.Minute
	if value, err := settings.Duration("values.duration"); err == nil {
		if want != value {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}
}

func TestDurationArray(t *testing.T) {
	settings := getSettings()

	// valid string array
	want := []time.Duration{12 * time.Second, 3 * time.Minute, 5 * time.Hour}
	if value, err := settings.DurationArray("duration-array"); err == nil {
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}
}

func TestDurationMap(t *testing.T) {
	settings := getSettings()

	// valid string map
	want := map[string]time.Duration{
		"seconds": 12 * time.Second,
		"minutes": 3 * time.Minute,
		"hours":   5 * time.Hour,
	}
	if value, err := settings.DurationMap("duration-map"); err == nil {
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}
}

func TestSize(t *testing.T) {
	settings := getSettings()
	var want int64 = 15 * 1024 * 1024 * 1024 * 1024
	if value, err := settings.Size("values.size"); err == nil {
		if want != value {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}
}

func TestSizeArray(t *testing.T) {
	settings := getSettings()

	want := []int64{
		15 * 1024,
		32 * 1024 * 1024,
		6 * 1024 * 1024 * 1024,
		2 * 1024 * 1024 * 1024 * 1024,
	}
	if value, err := settings.SizeArray("size-array"); err == nil {
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}
}

func TestSizeMap(t *testing.T) {
	settings := getSettings()

	want := map[string]int64{
		"kb": 15 * 1024,
		"mb": 32 * 1024 * 1024,
		"gb": 6 * 1024 * 1024 * 1024,
		"tb": 2 * 1024 * 1024 * 1024 * 1024,
	}
	if value, err := settings.SizeMap("size-map"); err == nil {
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}
}

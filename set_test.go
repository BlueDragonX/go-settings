package settings

import (
	"reflect"
	"testing"
)

func TestSet(t *testing.T) {
	var key string
	var want, value interface{}
	var settings *Settings

	// set a value in an empty settings object
	key = "key"
	want = "value"
	settings = &Settings{}
	settings.Set(key, want)
	value, _ = settings.String(key)
	if want != value {
		t.Errorf("%s != %s", want, value)
	}

	// set top level value
	settings = getSettings()
	key = "new-value"
	want = "Hello! I'm new here."
	settings.Set(key, want)
	value, _ = settings.String(key)
	if want != value {
		t.Errorf("%s != %s", want, value)
	}

	// set map value
	key = "mapping.a"
	want = "eh"
	settings.Set(key, want)
	value, _ = settings.String(key)
	if want != value {
		t.Errorf("%s != %s", want, value)
	}

	// set array value
	key = "integer-array.1"
	want = 3
	settings.Set(key, want)
	value, _ = settings.Int(key)
	if want != value {
		t.Errorf("%s != %s", want, value)
	}

	// set value to new map
	key = "new.map.here"
	want = "there"
	settings.Set(key, want)
	value, _ = settings.String(key)
	if want != value {
		t.Errorf("%s != %s", want, value)
	}

	// set non-interface map
	key = "new.map.there"
	want = map[interface{}]interface{}{"a": "aye", "b": "bee"}
	value = map[string]string{"a": "aye", "b": "bee"}
	settings.Set(key, value)
	value, _ = settings.Raw(key)
	if !reflect.DeepEqual(want, value) {
		t.Errorf("%v != %v", want, value)
	}

	// set non-interface array
	key = "new.array.there"
	want = []interface{}{"one", "two"}
	value = []string{"one", "two"}
	settings.Set(key, value)
	value, _ = settings.Raw(key)
	if !reflect.DeepEqual(want, value) {
		t.Errorf("%v != %v", want, value)
	}

	// set non-interface struct
	key = "new.struct.there"
	want = map[interface{}]interface{}{"A": "aye", "B": "bee"}
	value = struct {
		A string
		B string
		c string
	}{
		A: "aye",
		B: "bee",
		c: "see",
	}
	settings.Set(key, value)
	value, _ = settings.Raw(key)
	if !reflect.DeepEqual(want, value) {
		t.Errorf("%v != %v", want, value)
	}
}

func TestAppend(t *testing.T) {
	var err error
	var key string
	var want, value []string
	var settings *Settings

	// add to new root array in empty settings
	settings = &Settings{}
	key = "new-array"
	want = []string{"a"}
	err = settings.Append(key, "a")
	if err == nil {
		value, _ = settings.StringArray(key)
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}

	// add to new root array
	settings = getSettings()
	key = "new-array"
	want = []string{"a"}
	err = settings.Append(key, "a")
	if err == nil {
		value, _ = settings.StringArray(key)
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}

	// add to existing array with elements in it
	key = "new-array"
	want = []string{"a", "b"}
	err = settings.Append(key, "b")
	if err == nil {
		value, _ = settings.StringArray(key)
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}

	// add to existing root array
	key = "string-array"
	want = []string{"one", "two", "three"}
	err = settings.Append(key, "three")
	if err == nil {
		value, _ = settings.StringArray(key)
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}

	// add to new array in map
	key = "values.array"
	want = []string{"one"}
	err = settings.Append(key, "one")
	if err == nil {
		value, _ = settings.StringArray(key)
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}

	// add to existing array in map
	key = "values.array"
	want = []string{"one", "two"}
	err = settings.Append(key, "two")
	if err == nil {
		value, _ = settings.StringArray(key)
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}

	// add to new array in array
	key = "mixed-array.0"
	want = []string{"one"}
	err = settings.Append(key, "one")
	if err == nil {
		value, _ = settings.StringArray(key)
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}

	// add to existing array in array
	key = "mixed-array.0"
	want = []string{"one", "two"}
	err = settings.Append(key, "two")
	if err == nil {
		value, _ = settings.StringArray(key)
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	var err error
	var key string
	settings := getSettings()

	// delete item in root
	key = "value"
	settings.Set(key, "one")
	if err = settings.Delete(key); err == nil {
		if value, err := settings.Raw(key); err != KeyError {
			t.Errorf("%s not deleted: %v (%s)\n", key, value, err)
		}
	} else {
		t.Error(err)
	}

	// delete item in root map
	key = "values.bool"
	if err = settings.Delete(key); err == nil {
		if value, err := settings.Raw(key); err != KeyError {
			t.Errorf("%s not deleted: %v (%s)\n", key, value, err)
		}
	} else {
		t.Error(err)
	}

	// delete item in root array
	key = "string-array.one"
	if err = settings.Delete(key); err == nil {
		if value, err := settings.Raw(key); err != KeyError {
			t.Errorf("%s not deleted: %v (%s)\n", key, value, err)
		}
	} else {
		t.Error(err)
	}

	// delete item in nested map
	key = "settings-array.1.value"
	if err = settings.Delete(key); err == nil {
		if value, err := settings.Raw(key); err != KeyError {
			t.Errorf("%s not deleted: %v (%s)\n", key, value, err)
		}
	} else {
		t.Error(err)
	}

	// delete item in nested array
	want := []string{"two"}
	if err = settings.Delete("nested.array.0"); err == nil {
		value, _ := settings.StringArray("nested.array")
		if !reflect.DeepEqual(want, value) {
			t.Errorf("%s not deleted: %v\n", key, value)
		}
	} else {
		t.Error(err)
	}

	// delete a non-existent key
	key = "not-appearing-in-this-film"
	if err = settings.Delete(key); err != nil {
		t.Error(err)
	}
}

package settings

import "testing"

func isArrayEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for n := range a {
		if a[n] != b[n] {
			return false
		}
	}
	return true
}

func TestSet(t *testing.T) {
	var key string
	var want, value interface{}
	settings, _ := Parse([]byte(input))

	// set top level value
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
}

func TestAppend(t *testing.T) {
	var err error
	var key string
	var want, value []string
	settings, _ := Parse([]byte(input))

	// add to new root array
	key = "new-array"
	want = []string{"a"}
	err = settings.Append(key, "a")
	if err == nil {
		value, _ = settings.StringArray(key)
		if !isArrayEqual(want, value) {
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
		if !isArrayEqual(want, value) {
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
		if !isArrayEqual(want, value) {
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
		if !isArrayEqual(want, value) {
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
		if !isArrayEqual(want, value) {
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
		if !isArrayEqual(want, value) {
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
		if !isArrayEqual(want, value) {
			t.Errorf("%v != %v", want, value)
		}
	} else {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	var err error
	var key string
	settings, _ := Parse([]byte(input))

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
		if !isArrayEqual(want, value) {
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

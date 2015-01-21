package settings

import (
	"reflect"
	"testing"
)

func TestRawDflt(t *testing.T) {
	settings := getSettings()
	dflt := "nope"
	have := settings.RawDflt("nope", dflt)
	if !reflect.DeepEqual(dflt, have) {
		t.Errorf("%v != %v", dflt, have)
	}
}

func TestObjectDflt(t *testing.T) {
	settings := getSettings()
	var dflt *Settings = nil
	have := settings.ObjectDflt("nope", dflt)
	if !reflect.DeepEqual(dflt, have) {
		t.Errorf("%v != %v", dflt, have)
	}
}

func TestObjectArrayDflt(t *testing.T) {
	settings := getSettings()
	dflt := []*Settings{}
	have := settings.ObjectArrayDflt("nope", dflt)
	if !reflect.DeepEqual(dflt, have) {
		t.Errorf("%v != %v", dflt, have)
	}
}

func TestObjectMapDflt(t *testing.T) {
	settings := getSettings()
	dflt := map[string]*Settings{}
	have := settings.ObjectMapDflt("nope", dflt)
	if !reflect.DeepEqual(dflt, have) {
		t.Errorf("%v != %v", dflt, have)
	}
}

func TestStringDflt(t *testing.T) {
	settings := getSettings()
	dflt := "hello"
	have := settings.StringDflt("nope", dflt)
	if !reflect.DeepEqual(dflt, have) {
		t.Errorf("%v != %v", dflt, have)
	}
}

func TestStringArrayDflt(t *testing.T) {
	settings := getSettings()
	dflt := []string{"one", "two"}
	have := settings.StringArrayDflt("nope", dflt)
	if !reflect.DeepEqual(dflt, have) {
		t.Errorf("%v != %v", dflt, have)
	}
}

func TestStringMapDflt(t *testing.T) {
	settings := getSettings()
	dflt := map[string]string{"a": "aye", "b": "bee"}
	have := settings.StringMapDflt("nope", dflt)
	if !reflect.DeepEqual(dflt, have) {
		t.Errorf("%v != %v", dflt, have)
	}
}

func TestIntDflt(t *testing.T) {
	settings := getSettings()
	dflt := 12
	have := settings.IntDflt("nope", dflt)
	if !reflect.DeepEqual(dflt, have) {
		t.Errorf("%v != %v", dflt, have)
	}
}

func TestIntArrayDflt(t *testing.T) {
	settings := getSettings()
	dflt := []int{1, 2, 3}
	have := settings.IntArrayDflt("nope", dflt)
	if !reflect.DeepEqual(dflt, have) {
		t.Errorf("%v != %v", dflt, have)
	}
}

func TestIntMapDflt(t *testing.T) {
	settings := getSettings()
	dflt := map[string]int{"one": 1, "two": 2}
	have := settings.IntMapDflt("nope", dflt)
	if !reflect.DeepEqual(dflt, have) {
		t.Errorf("%v != %v", dflt, have)
	}
}

func TestFloatDflt(t *testing.T) {
	settings := getSettings()
	dflt := 15.6
	have := settings.FloatDflt("nope", dflt)
	if !reflect.DeepEqual(dflt, have) {
		t.Errorf("%v != %v", dflt, have)
	}
}

func TestFloatArrayDflt(t *testing.T) {
	settings := getSettings()
	dflt := []float64{1.3, 2.9, 16.8}
	have := settings.FloatArrayDflt("nope", dflt)
	if !reflect.DeepEqual(dflt, have) {
		t.Errorf("%v != %v", dflt, have)
	}
}

func TestFloatMapDflt(t *testing.T) {
	settings := getSettings()
	dflt := map[string]float64{"one": 1.2, "two": 2.3}
	have := settings.FloatMapDflt("nope", dflt)
	if !reflect.DeepEqual(dflt, have) {
		t.Errorf("%v != %v", dflt, have)
	}
}

func TestBoolDflt(t *testing.T) {
	settings := getSettings()
	dflt := true
	have := settings.BoolDflt("nope", dflt)
	if !reflect.DeepEqual(dflt, have) {
		t.Errorf("%v != %v", dflt, have)
	}
}

func TestBoolArrayDflt(t *testing.T) {
	settings := getSettings()
	dflt := []bool{true, false, false}
	have := settings.BoolArrayDflt("nope", dflt)
	if !reflect.DeepEqual(dflt, have) {
		t.Errorf("%v != %v", dflt, have)
	}
}

func TestBoolMapDflt(t *testing.T) {
	settings := getSettings()
	dflt := map[string]bool{"yes": true, "no": false}
	have := settings.BoolMapDflt("nope", dflt)
	if !reflect.DeepEqual(dflt, have) {
		t.Errorf("%v != %v", dflt, have)
	}
}

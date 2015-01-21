package settings

// Get a value from the settings object. Return `dflt` if an error occurs.
func (s *Settings) RawDflt(key string, dflt interface{}) interface{} {
	if value, err := s.Raw(key); err == nil {
		return value
	} else {
		return dflt
	}
}

// Get a settings object. Return `dflt` if an error occurs.
func (s *Settings) ObjectDflt(key string, dflt *Settings) *Settings {
	if value, err := s.Object(key); err == nil {
		return value
	} else {
		return dflt
	}
}

// Get an array of settings objects. Return `dflt` if an error occurs.
func (s *Settings) ObjectArrayDflt(key string, dflt []*Settings) []*Settings {
	if value, err := s.ObjectArray(key); err == nil {
		return value
	} else {
		return dflt
	}
}

// Get a map of settings objects. Return `dflt` if an error occurs.
func (s *Settings) ObjectMapDflt(key string, dflt map[string]*Settings) map[string]*Settings {
	if value, err := s.ObjectMap(key); err == nil {
		return value
	} else {
		return dflt
	}
}

// Get a string value. Return `dflt` if an error occurs.
func (s *Settings) StringDflt(key string, dflt string) string {
	if value, err := s.String(key); err == nil {
		return value
	} else {
		return dflt
	}
}

// Get an array of string values. Return `dflt` if an error occurs.
func (s *Settings) StringArrayDflt(key string, dflt []string) []string {
	if value, err := s.StringArray(key); err == nil {
		return value
	} else {
		return dflt
	}
}

// Get a map of string values. Return `dflt` if an error occurs.
func (s *Settings) StringMapDflt(key string, dflt map[string]string) map[string]string {
	if value, err := s.StringMap(key); err == nil {
		return value
	} else {
		return dflt
	}
}

// Get an integer value. Return `dflt` if an error occurs.
func (s *Settings) IntDflt(key string, dflt int) int {
	if value, err := s.Int(key); err == nil {
		return value
	} else {
		return dflt
	}
}

// Get an array of integer values. Return `dflt` if an error occurs.
func (s *Settings) IntArrayDflt(key string, dflt []int) []int {
	if value, err := s.IntArray(key); err == nil {
		return value
	} else {
		return dflt
	}
}

// Get a map of integer values. Return `dflt` if an error occurs.
func (s *Settings) IntMapDflt(key string, dflt map[string]int) map[string]int {
	if value, err := s.IntMap(key); err == nil {
		return value
	} else {
		return dflt
	}
}

// Get a float value. Return `dflt` if an error occurs.
func (s *Settings) FloatDflt(key string, dflt float64) float64 {
	if value, err := s.Float(key); err == nil {
		return value
	} else {
		return dflt
	}
}

// Get an array of float values. Return `dflt` if an error occurs.
func (s *Settings) FloatArrayDflt(key string, dflt []float64) []float64 {
	if value, err := s.FloatArray(key); err == nil {
		return value
	} else {
		return dflt
	}
}

// Get a map of float values. Return `dflt` if an error occurs.
func (s *Settings) FloatMapDflt(key string, dflt map[string]float64) map[string]float64 {
	if value, err := s.FloatMap(key); err == nil {
		return value
	} else {
		return dflt
	}
}

// Get a bool value. Return `dflt` if an error occurs.
func (s *Settings) BoolDflt(key string, dflt bool) bool {
	if value, err := s.Bool(key); err == nil {
		return value
	} else {
		return dflt
	}
}

// Get an array of bool values. Return `dflt` if an error occurs.
func (s *Settings) BoolArrayDflt(key string, dflt []bool) []bool {
	if value, err := s.BoolArray(key); err == nil {
		return value
	} else {
		return dflt
	}
}

// Get a map of bool values. Return `dflt` if an error occurs.
func (s *Settings) BoolMapDflt(key string, dflt map[string]bool) map[string]bool {
	if value, err := s.BoolMap(key); err == nil {
		return value
	} else {
		return dflt
	}
}

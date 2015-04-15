package settings

import (
	"strconv"
	"strings"
)

// Parse a size string into a number of bytes. Valid suffixes are b, k, kb, m,
// mb, g, gb, t, and tb. Case is ignored. A lack of suffix indicates bytes.
func ParseSize(s string) (int64, error) {
	type size struct {
		symbol string
		factor int64
	}

	sizes := []size{
		{"kb", 1024},
		{"mb", 1048576},
		{"gb", 1073741824},
		{"tb", 1099511627776},
		{"b", 1},
		{"k", 1024},
		{"m", 1048576},
		{"g", 1073741824},
		{"t", 1099511627776},
	}

	s = strings.TrimSpace(strings.ToLower(s))
	var factor int64 = 1
	for _, size := range sizes {
		pos := len(s) - len(size.symbol)
		if s[pos:] == size.symbol {
			s = strings.TrimSpace(s[:pos])
			factor = size.factor
			break
		}
	}

	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return n * factor, nil
	} else {
		return 0, err
	}
}

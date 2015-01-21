package settings

import "errors"

var IndexError error = errors.New("invalid index")
var KeyError error = errors.New("key not found")
var ObjectError error = errors.New("invalid object")
var RangeError error = errors.New("index out of range")
var TypeError error = errors.New("invalid type conversion")

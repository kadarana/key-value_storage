package storage

import (
	"errors"
	"math"
)
// used to set values to a structure 
func newValue(val any) (Value, error) {
	valueType := getType(val)
	if valueType != kindUndefind {
		return Value{
			Val:       val,
			ValueType: valueType,
		}, nil
	}
	return Value{}, errors.New("Undefined ValueType")
}

// used to set const to a value's type
func getType(val any) Kind {
	switch val.(type) {
	case int:
		return kindInt
	case float64:
		if isFloatInt(val) {
			return kindInt
		}
		return kindUndefind
	case string:
		return kindString
	default:
		return kindUndefind
	}
}


func isFloatInt(num any) bool {
	return num.(float64) == math.Trunc(num.(float64))
}

// subfunction for HGET
func (r *Storage) hget(key string, field string) (Value, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	res, ok := r.innerMap[key][field]
	if !ok {
		return Value{}, false
	}
	return res, true
}


// using for function LPOP
func convertIndex(index int, length int) int {
	if index < 0 {
		index = length + index
	}

	return index
}

// using for function RPOP
func normalizeIndex(index, length int) int {
	if index < 0 {
		index = length + index
		if index < 0 {
			index = 0
		}
	}
	return index
}

// using for function RPOP
func reverse(slice []any) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}
package conver

import (
	"fmt"
	"strconv"
	"strings"
)

// String : Conver "val" to a String
func String(val interface{}) (string, error) {
	switch ret := val.(type) {
	case string:
		return ret, nil
	case []byte:
		return string(ret), nil
	default:
		str := fmt.Sprintf("%+v", val)
		if val == nil || len(str) == 0 {
			return "", fmt.Errorf("conver.String(), the %+v is empty", val)
		}
		return str, nil
	}
}

// StringMust : Must Conver "val" to a String
func StringMust(val interface{}, def ...string) string {
	ret, err := String(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Bool : Conver "val" to a Bool
func Bool(val interface{}) (bool, error) {
	if val == nil {
		return false, nil
	}
	switch ret := val.(type) {
	case bool:
		return ret, nil
	case int, int8, int16, int32, int64, float32, float64, uint, uint8, uint16, uint32, uint64:
		return ret != 0, nil
	case []byte:
		return stringToBool(string(ret))
	case string:
		return stringToBool(ret)
	default:
		return false, converError(val, "bool")
	}
}

// BoolMust : Must Conver "val" to a Bool
func BoolMust(val interface{}, def ...bool) bool {
	ret, err := Bool(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Bytes : Conver "val" to []byte
func Bytes(val interface{}) ([]byte, error) {
	switch ret := val.(type) {
	case []byte:
		return ret, nil
	default:
		str, err := String(val)
		return []byte(str), err
	}
}

// BytesMust : Must Conver "val" to []byte
func BytesMust(val interface{}, def ...[]byte) []byte {
	ret, err := Bytes(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Float32 : Conver "val" to a Float32
func Float32(val interface{}) (float32, error) {
	switch ret := val.(type) {
	case float32:
		return ret, nil
	case int:
		return float32(ret), nil
	case int8:
		return float32(ret), nil
	case int16:
		return float32(ret), nil
	case int32:
		return float32(ret), nil
	case int64:
		return float32(ret), nil
	case uint:
		return float32(ret), nil
	case uint8:
		return float32(ret), nil
	case uint16:
		return float32(ret), nil
	case uint32:
		return float32(ret), nil
	case uint64:
		return float32(ret), nil
	case float64:
		return float32(ret), nil
	case bool:
		if ret {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		str := strings.Replace(strings.TrimSpace(StringMust(val)), " ", "", -1)
		f, err := strconv.ParseFloat(str, 32)
		return float32(f), err
	}
}

// Float32Must : Must Conver "val" to Float32
func Float32Must(val interface{}, def ...float32) float32 {
	ret, err := Float32(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Float64 : Conver "val" to a Float64
func Float64(val interface{}) (float64, error) {
	switch ret := val.(type) {
	case float64:
		return ret, nil
	case int:
		return float64(ret), nil
	case int8:
		return float64(ret), nil
	case int16:
		return float64(ret), nil
	case int32:
		return float64(ret), nil
	case int64:
		return float64(ret), nil
	case uint:
		return float64(ret), nil
	case uint8:
		return float64(ret), nil
	case uint16:
		return float64(ret), nil
	case uint32:
		return float64(ret), nil
	case uint64:
		return float64(ret), nil
	case float32:
		return float64(ret), nil
	case bool:
		if ret {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		str := strings.Replace(strings.TrimSpace(StringMust(val)), " ", "", -1)
		return strconv.ParseFloat(str, 64)
	}
}

// Float64Must : Must Conver "val" to Float64
func Float64Must(val interface{}, def ...float64) float64 {
	ret, err := Float64(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Int : Conver "val" to a rounded Int
func Int(val interface{}) (int, error) {
	f, err := Float64(val)
	if err != nil {
		return 0, err
	}
	str := strconv.FormatFloat(f, 'f', 0, 64)
	i, err := strconv.ParseInt(str, 10, 0)
	return int(i), err
}

// IntMust : Must Conver "val" to a rounded Int
func IntMust(val interface{}, def ...int) int {
	ret, err := Int(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Int32 : Conver "val" to a rounded Int32
func Int32(val interface{}) (int32, error) {
	f, err := Float64(val)
	if err != nil {
		return 0, err
	}
	str := strconv.FormatFloat(f, 'f', 0, 64)
	i, err := strconv.ParseInt(str, 10, 32)
	return int32(i), err
}

// Int32Must : Must Conver "val" to a rounded Int32
func Int32Must(val interface{}, def ...int32) int32 {
	ret, err := Int32(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Int64 : Conver "val" to a rounded Int64
func Int64(val interface{}) (int64, error) {
	f, err := Float64(val)
	if err != nil {
		return 0, err
	}
	str := strconv.FormatFloat(f, 'f', 0, 64)
	return strconv.ParseInt(str, 10, 64)
}

// Int64Must : Must Conver "val" to a rounded Int64
func Int64Must(val interface{}, def ...int64) int64 {
	ret, err := Int64(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

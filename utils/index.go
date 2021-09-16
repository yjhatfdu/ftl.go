package utils
func Truly(v interface{}) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val != ""
	case int8, int16, int, int32, int64, uint8, uint16, uint32, uint64, float32, float64:
		return val != 0
	default:
		return true
	}
}


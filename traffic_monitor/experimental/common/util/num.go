package util

// ToNumeric returns a float for any numeric type, and false if the interface does not hold a numeric type.
// This allows converting unknown numeric types (for example, from JSON) in a single line
// TODO try to parse string stats as numbers?
func ToNumeric(v interface{}) (float64, bool) {
	switch i := v.(type) {
	case uint8:
		return float64(i), true
	case uint16:
		return float64(i), true
	case uint32:
		return float64(i), true
	case uint64:
		return float64(i), true
	case int8:
		return float64(i), true
	case int16:
		return float64(i), true
	case int32:
		return float64(i), true
	case int64:
		return float64(i), true
	case float32:
		return float64(i), true
	case float64:
		return i, true
	case int:
		return float64(i), true
	case uint:
		return float64(i), true
	default:
		return 0.0, false
	}
}

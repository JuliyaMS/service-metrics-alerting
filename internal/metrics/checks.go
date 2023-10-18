package metrics

import "strconv"

func CheckType(value string) bool {
	types := []string{"gauge", "counter"}
	for _, tp := range types {
		if tp == value {
			return true
		}
	}
	return false
}

func CheckDigit(t string, val string) bool {
	if t == "gauge" {
		if _, err := strconv.ParseFloat(val, 64); err == nil {
			return true
		}
	}
	if t == "counter" {
		if _, err := strconv.ParseInt(val, 10, 64); err == nil {
			return true
		}
	}
	return false
}

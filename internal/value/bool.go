package value

import "github.com/Gulliby-development/cfgquill/internal/base"

// Bool interprets the Git boolean spellings and short decimal forms.
func Bool(v string) (value bool, ok bool) {
	switch base.Lower(v) {
	case "true", "yes", "on":
		return true, true
	case "", "false", "no", "off":
		return false, true
	}

	if len(v) == 0 || len(v) > 9 {
		return false, false
	}
	nonzero := false
	for i := 0; i < len(v); i++ {
		if v[i] < '0' || v[i] > '9' {
			return false, false
		}
		if v[i] != '0' {
			nonzero = true
		}
	}
	return nonzero, true
}

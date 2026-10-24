package problem

import "github.com/Gulliby-development/cfgquill/internal/base"

// FaultText returns the diagnostic text associated with a syntax code.
func FaultText(code, arg string) string {
	switch code {
	case "CQ001":
		return "bad section header"
	case "CQ002":
		return "bad variable name"
	case "CQ003":
		return "expected \"=\" after the variable name"
	case "CQ004":
		return "unterminated quote"
	case "CQ005":
		return "unknown escape " + base.Show(arg)
	case "CQ006":
		return "text after the end of the value"
	default:
		return ""
	}
}

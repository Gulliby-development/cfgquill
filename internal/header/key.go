package header

import "github.com/Gulliby-development/cfgquill/internal/base"

// Canonical validates a key and converts its section and variable to lowercase.
func Canonical(key string) (string, bool) {
	firstDot, lastDot := -1, -1
	for i := 0; i < len(key); i++ {
		if key[i] == '.' {
			if firstDot == -1 {
				firstDot = i
			}
			lastDot = i
		}
	}
	if firstDot <= 0 || lastDot == len(key)-1 {
		return "", false
	}

	section := key[:firstDot]
	variable := key[lastDot+1:]
	for _, r := range section {
		if !base.IsKeyChar(r) {
			return "", false
		}
	}

	firstVariableRune := true
	for _, r := range variable {
		if firstVariableRune {
			if !base.IsAlpha(r) {
				return "", false
			}
			firstVariableRune = false
		} else if !base.IsKeyChar(r) {
			return "", false
		}
	}
	if firstVariableRune {
		return "", false
	}

	var subsection string
	if lastDot > firstDot {
		subsection = key[firstDot+1 : lastDot]
		for _, r := range subsection {
			if r == '\n' {
				return "", false
			}
		}
	}

	canonical := base.Lower(section) + "."
	if lastDot > firstDot {
		canonical += subsection + "."
	}
	canonical += base.Lower(variable)
	return canonical, true
}

// Split returns the section and subsection portions of a stored name.
func Split(name string) (section, sub string, hasSub bool) {
	for i := 0; i < len(name); i++ {
		if name[i] == '.' {
			return name[:i], name[i+1:], true
		}
	}
	return name, "", false
}

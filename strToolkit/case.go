package strToolkit

import (
	"regexp"
	"strings"
)

var uppercaseAcronym = map[string]bool{
	"ID": true,
}

func IsUpperCase(r rune) bool {
	if r >= 'A' && r <= 'Z' {
		return true
	}
	return false
}

func IsLowerCase(r rune) bool {
	if r >= 'a' && r <= 'z' {
		return true
	}
	return false
}

// Converts a string to CamelCase
func toCamelInitCase(s string, initCase bool) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	capNext := initCase
	for _, v := range s {
		if v >= 'A' && v <= 'Z' {
			n += string(v)
		}
		if v >= '0' && v <= '9' {
			n += string(v)
		}
		if v >= 'a' && v <= 'z' {
			if capNext {
				n += strings.ToUpper(string(v))
			} else {
				n += string(v)
			}
		}
		if v == '_' || v == ' ' || v == '-' || v == '.' {
			capNext = true
		} else {
			capNext = false
		}
	}
	return n
}

// ToCamel converts a string to CamelCase
func ToCamel(s string) string {
	if uppercaseAcronym[s] {
		s = strings.ToLower(s)
	}
	return toCamelInitCase(s, true)
}

// ToLowerCamel converts a string to lowerCamelCase
func ToLowerCamel(s string) string {
	if s == "" {
		return s
	}
	if uppercaseAcronym[s] {
		s = strings.ToLower(s)
	}
	if r := rune(s[0]); r >= 'A' && r <= 'Z' {
		s = strings.ToLower(string(r)) + s[1:]
	}
	return toCamelInitCase(s, false)
}

// ToSnake converts a string to snake_case
func ToSnake(s string) string {

	return ToDelimited(s, '_')
}
func ToSnakeWithIgnore(s string, ignore uint8) string {

	return ToScreamingDelimited(s, '_', ignore, false)
}

// ToScreamingSnake converts a string to SCREAMING_SNAKE_CASE
func ToScreamingSnake(s string) string {
	return ToScreamingDelimited(s, '_', 0, true)
}

// ToKebab converts a string to kebab-case
func ToKebab(s string) string {
	return ToDelimited(s, '-')
}

// ToScreamingKebab converts a string to SCREAMING-KEBAB-CASE
func ToScreamingKebab(s string) string {
	return ToScreamingDelimited(s, '-', 0, true)
}

// ToDelimited converts a string to delimited.snake.case
// (in this case `delimiter = '.'`)
func ToDelimited(s string, delimiter uint8) string {
	return ToScreamingDelimited(s, delimiter, 0, false)
}

// ToScreamingDelimited converts a string to SCREAMING.DELIMITED.SNAKE.CASE
// (in this case `delimiter = '.'; screaming = true`)
// or delimited.snake.case
// (in this case `delimiter = '.'; screaming = false`)
func ToScreamingDelimited(s string, delimiter uint8, ignore uint8, screaming bool) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	for i, v := range s {
		// treat acronyms as words, eg for JSONData -> JSON is a whole word
		nextCaseIsChanged := false
		if i+1 < len(s) {
			next := s[i+1]
			vIsCap := v >= 'A' && v <= 'Z'
			vIsLow := v >= 'a' && v <= 'z'
			nextIsCap := next >= 'A' && next <= 'Z'
			nextIsLow := next >= 'a' && next <= 'z'
			if (vIsCap && nextIsLow) || (vIsLow && nextIsCap) {
				nextCaseIsChanged = true
			}
			if ignore > 0 && i-1 >= 0 && s[i-1] == ignore && nextCaseIsChanged {
				nextCaseIsChanged = false
			}
		}

		if i > 0 && n[len(n)-1] != delimiter && nextCaseIsChanged {
			// add underscore if next letter case type is changed
			if v >= 'A' && v <= 'Z' {
				n += string(delimiter) + string(v)
			} else if v >= 'a' && v <= 'z' {
				n += string(v) + string(delimiter)
			}
		} else if v == ' ' || v == '_' || v == '-' {
			// replace spaces/underscores with delimiters
			if uint8(v) == ignore {
				n += string(v)
			} else {
				n += string(delimiter)
			}
		} else {
			n = n + string(v)
		}
	}

	if screaming {
		n = strings.ToUpper(n)
	} else {
		n = strings.ToLower(n)
	}
	return n
}

var numberSequence = regexp.MustCompile(`([a-zA-Z])(\d+)([a-zA-Z]?)`)
var numberReplacement = []byte(`$1 $2 $3`)

func addWordBoundariesToNumbers(s string) string {
	b := []byte(s)
	b = numberSequence.ReplaceAll(b, numberReplacement)
	return string(b)
}

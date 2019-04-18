package strToolkit

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

func ToSnakeCase(s string) string {
	out := []rune{}
	for index, r := range s {
		if index == 0 {
			out = append(out, ToLowerCase(r))
			continue
		}

		if IsUpperCase(r) && index != 0 {
			out = append(out, '_', ToLowerCase(r))
			continue
		}
		out = append(out, r)
	}
	return string(out)
}

func ToCamelCase(s string) string {
	out := []rune{}
	for index, r := range s {
		if r == '_' {
			continue
		}
		if index == 0 {
			out = append(out, ToUpperCase(r))
			continue
		}

		if index > 0 && s[index-1] == '_' {
			out = append(out, ToUpperCase(r))
			continue
		}

		out = append(out, r)
	}
	return string(out)
}

func ToLowerCase(r rune) rune {
	dx := 'A' - 'a'
	if IsUpperCase(r) {
		return r - dx
	}
	return r
}
func ToUpperCase(r rune) rune {
	dx := 'A' - 'a'
	if IsLowerCase(r) {
		return r + dx
	}
	return r
}

func ToLower(s string) string {
	out := []rune{}
	for _, r := range s {
		out = append(out, ToLowerCase(r))
	}
	return string(out)
}

func ToUpper(s string) string {
	out := []rune{}
	for _, r := range s {
		out = append(out, ToUpperCase(r))
	}
	return string(out)
}

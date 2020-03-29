package strToolkit

func SlicifyStr(seed string, length int) []string {
	slice := make([]string, length)
	for i := 0; i < length; i++ {
		slice[i] = seed
	}
	return slice
}

func SliceContains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

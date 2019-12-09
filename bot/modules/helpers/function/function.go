package function

func Contains(key []string, str string) bool {
	if key != nil && str != "" {
		for _, val := range key {
			if val == str {
				return true
			}
		}
		return false
	}
	return false
}

package function

func Contains(s []string, e string) bool {
	if s != nil && e != "" {
		for _, a := range s {
			if a == e {
				return true
			}
		}
		return false
	}
	return false
}

package slice

func ContainsString(strings []string, check string) bool {
	for _, s := range strings {
		if s == check {
			return true
		}
	}
	return false
}

func ContainsInt64(ints []int64, check int64) bool {
	for _, s := range ints {
		if s == check {
			return true
		}
	}
	return false
}

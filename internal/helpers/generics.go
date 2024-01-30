package helpers

func Filter[T any](ss []T, filter func(T) bool) (retVal []T) {
	for _, s := range ss {
		if filter(s) {
			retVal = append(retVal, s)
		}
	}

	return
}

func FindFirst[T any](ss []T, filter func(T) bool) (retVal T) {
	for _, s := range ss {
		if filter(s) {
			retVal = s
			break
		}
	}

	return
}

func AnyMatch[T any](ss []T, filter func(T) bool) bool {
	for _, s := range ss {
		if filter(s) {
			return true
		}
	}

	return false
}

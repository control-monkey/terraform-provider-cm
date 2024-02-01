package helpers

func Filter[T any](ss []T, filter func(T) bool) (retVal []T) {
	for _, s := range ss {
		if filter(s) {
			retVal = append(retVal, s)
		}
	}

	return
}

func Map[T any, S any](ss []T, f func(T) S) (retVal []S) {
	for _, s := range ss {
		m := f(s)
		retVal = append(retVal, m)
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

func IsUnique[T int | string](es []T) bool {
	alreadyAppeared := make(map[T]bool)

	for _, e := range es {
		if _, value := alreadyAppeared[e]; !value {
			alreadyAppeared[e] = true
		} else {
			return false
		}
	}

	return true
}

func FindDuplicates[T int | string](es []T, stopAfterFirst bool) []T {
	retVal := make([]T, 0)
	keys := make(map[T]bool)

	for _, e := range es {
		if _, value := keys[e]; !value {
			keys[e] = true
		} else {
			retVal = append(retVal, e)
			if stopAfterFirst {
				break
			}
		}
	}

	return retVal
}

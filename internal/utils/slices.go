package utils

func Contains[T comparable](slice []T, a T) bool {
	for _, item := range slice {
		if item == a {
			return true
		}
	}
	return false
}

func Index[T comparable](slice []T, a T) int {
	for idx, item := range slice {
		if item == a {
			return idx
		}
	}
	return -1
}

func RemoveIdx[T comparable](slice []T, idx int) []T {
	if idx >= len(slice) {
		return slice[0:idx]
	}
	return append(slice[0:idx], slice[idx+1])
}

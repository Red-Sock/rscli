package slices

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
	return append(slice[0:idx], slice[idx+1:]...)
}

func Exclude[T comparable](slice []T, elems ...T) []T {
	for _, itemToExclude := range elems {
		for idx, item := range slice {
			if item == itemToExclude {
				slice = RemoveIdx(slice, idx)
				break
			}
		}
	}
	return slice
}

func InsertSlice[T any](src []T, insert []T, idx int) []T {
	afterSlice := make([]T, len(src[idx:]))
	copy(afterSlice, src[idx:])

	beforeSlice := make([]T, len(src[:idx]))
	copy(beforeSlice, src[:idx])

	return append(append(beforeSlice, insert...), afterSlice...)
}

func RemovePart[T any](src []T, start, end int) []T {
	end++
	out := make([]T, len(src[:start])+len(src[end:]))
	copy(out[:start], src[:start])
	copy(out[start:], src[end:])
	return out
}

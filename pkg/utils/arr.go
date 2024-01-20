package utils

func MapArr[T comparable, V comparable](arr []T, callback func(v T) V) []V {
	newArr := make([]V, len(arr))
	for i, v := range arr {
		newArr[i] = callback(v)
	}
	return newArr
}

func FilterArr[T comparable](arr []T, callback func(v T) bool) []T {
	newArr := make([]T, 0)
	for _, v := range arr {
		if callback(v) {
			newArr = append(newArr, v)
		}
	}
	return newArr
}

package utils

func Insert[T any](orig []T, index int, value T) []T {
	if index < 0 {
		panic("index out of range")
	}

	if index >= len(orig) {
		return append(orig, value)
	}

	orig = append(orig[:index+1], orig[index:]...)
	orig[index] = value

	return orig
}

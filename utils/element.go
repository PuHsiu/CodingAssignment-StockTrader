package utils

func Min[T int | int64 | int32 | int16 | int8](this, that T) T {
	if this >= that {
		return that
	}

	return this
}

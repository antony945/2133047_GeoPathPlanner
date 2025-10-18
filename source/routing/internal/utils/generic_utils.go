package utils

func GetOrDefault[T any](m map[string]any, key string, def T) T {
	if m == nil {
		return def
	}

	if v, ok := m[key]; ok {
		if cast, ok2 := v.(T); ok2 {
			return cast
		}
	}
	return def
}

func GetOrNil[T any](m map[string]any, key string) T {
	var zero T

	if m == nil {
		return zero
	}

	if v, ok := m[key]; ok {
		if cast, ok2 := v.(T); ok2 {
			return cast
		}
	}
	return zero
}
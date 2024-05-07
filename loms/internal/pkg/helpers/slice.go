// Package helpers Хелперы
package helpers

// Map Применяет отображение f к элементам слайса
func Map[T any, V any](in []T, f func(T) V) []V {
	res := make([]V, len(in))
	for i, t := range in {
		res[i] = f(t)
	}
	return res
}

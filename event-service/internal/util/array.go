package util

func Truncate[T any](arr []T, n int) []T {
	limit := n
	if len(arr) < limit {
		limit = len(arr)
	}
	return arr[:limit]
}

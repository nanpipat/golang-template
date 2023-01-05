package helper

func ArrayMap[T any, R any](raws []T, cb func(r T) R) []R {
	items := make([]R, 0)

	for _, r := range raws {
		items = append(items, cb(r))
	}

	return items
}

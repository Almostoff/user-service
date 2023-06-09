package utils

import "math"

func CalculateOffset(limit, page int64) int64 {
	page = int64(math.Max(1, float64(page)))
	return limit * (page - 1)
}

func CalculatePages(limit, _len int64) int64 {

	if _len < limit {
		return 1
	}

	var (
		noRoundedDiv float64 = float64(_len) / float64(limit)
		roundedDiv   int64   = _len / limit
		needAdd1     bool    = !(int64(noRoundedDiv*1000) == (roundedDiv * 1000))
	)

	if needAdd1 {
		return roundedDiv + 1
	}
	return roundedDiv
}

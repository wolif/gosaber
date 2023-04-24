package morm

var pageSizeMax = 1000

func SetPageSizeMax(n int) {
	if n <= 0 {
		return
	}
	pageSizeMax = n
}

func minInt(a, b int) int {
	if a > b {
		return b
	}

	return a
}

func resolvePagination(pagination ...int) (page int, pageSize int) {
	if len(pagination) >= 1 {
		page = pagination[0]
	}
	if len(pagination) >= 2 {
		pageSize = pagination[1]
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	pageSize = minInt(pageSize, pageSizeMax)
	return
}

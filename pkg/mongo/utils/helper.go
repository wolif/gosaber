package utils

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
	return
}

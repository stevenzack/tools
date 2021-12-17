package numToolkit

func PageNum(total int64, pageSize int) int {
	mod := int(total) % pageSize
	page := int(total) / pageSize
	if mod == 0 {
		return page
	}
	return page + 1
}

func ValidatePage(pageSize, page int, total int64) bool {
	if page == 0 {
		page = 1
	}
	if page < 1 {
		return false
	}
	if pageSize*(page-1) > int(total) {
		return false
	}
	return true
}

func OffsetBy(pageSize, page int) int64 {
	if page <= 0 {
		return 0
	}
	return int64(pageSize * (page - 1))
}

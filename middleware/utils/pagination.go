package utils

import (
	"math"
)

type PaginationParams struct {
	Page   int
	Limit  int
	Offset int
}

// ParsePagination parses and validates pagination parameters
func ParsePagination(page, limit int) PaginationParams {
	// Default values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	// Max limit to prevent abuse
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	return PaginationParams{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

// CalculateTotalPages calculates total pages from total items and limit
func CalculateTotalPages(total int64, limit int) int {
	if limit <= 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(limit)))
}

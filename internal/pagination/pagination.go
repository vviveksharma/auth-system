// pagination/pagination.go
package pagination

import (
	"math"
	"strconv"
)

type PaginationRequest struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

type PaginationMeta struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	TotalPages int  `json:"total_pages"`
	TotalItems int  `json:"total_items"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

type PaginatedResponse[T any] struct {
	Data       []T            `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

func ParsePaginationParams(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func PaginateSlice[T any](items []T, page, pageSize int) PaginatedResponse[T] {
	page, pageSize = ParsePaginationParams(page, pageSize)

	totalItems := len(items)
	totalPages := calculateTotalPages(totalItems, pageSize)

	// Calculate start and end indices
	start := (page - 1) * pageSize
	if start > totalItems {
		start = totalItems
	}

	end := start + pageSize
	if end > totalItems {
		end = totalItems
	}

	// Get paginated data
	var data []T
	if start < totalItems {
		data = items[start:end]
	} else {
		data = []T{}
	}

	// Build pagination metadata
	pagination := PaginationMeta{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalItems: totalItems,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	return PaginatedResponse[T]{
		Data:       data,
		Pagination: pagination,
	}
}

func PaginateFromQuery[T any](items []T, pageStr, pageSizeStr string) PaginatedResponse[T] {
	page, pageSize := parseQueryParams(pageStr, pageSizeStr)
	return PaginateSlice(items, page, pageSize)
}

func calculateTotalPages(totalItems, pageSize int) int {
	if totalItems == 0 || pageSize == 0 {
		return 0
	}
	return int(math.Ceil(float64(totalItems) / float64(pageSize)))
}


func parseQueryParams(pageStr, pageSizeStr string) (int, int) {
	page := 1
	pageSize := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
			if pageSize > 100 {
				pageSize = 100
			}
		}
	}

	return page, pageSize
}

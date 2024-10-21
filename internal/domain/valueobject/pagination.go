package valueobject

import "math"

const (
	PaginationDefaultPage     = 1
	PaginationDefaultPageSize = 10
)

// Pagination is the pagination configuration.
// Recommended to use `NewPagination` to create a new `Pagination` struct, because it sets default values.
type Pagination struct {
	page     int
	pageSize int
}

func (p Pagination) Limit() int {
	return p.pageSize
}

func (p Pagination) Offset() int {
	// handle negative values
	if p.page <= 0 || p.pageSize <= 0 {
		return 0
	}
	return (p.page - 1) * p.pageSize
}

func (p Pagination) CurrentPage() int {
	return p.page
}

// NewPagination creates a new `Pagination` struct with default values if none are provided.
func NewPagination(page, pageSize int) Pagination {
	if page <= 0 {
		page = PaginationDefaultPage
	}
	if pageSize <= 0 {
		pageSize = PaginationDefaultPageSize
	}

	return Pagination{
		page:     page,
		pageSize: pageSize,
	}
}

func DefaultPagination() Pagination {
	return NewPagination(PaginationDefaultPage, PaginationDefaultPageSize)
}

// PaginationMetadata represents the metadata for paginated responses.
type PaginationMetadata struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	FirstPage    int `json:"first_page"`
	LastPage     int `json:"last_page"`
	TotalRecords int `json:"total_records"`
}

// CalculatePaginationMetadata calculates the pagination metadata based on the total number of records, the current page, and the page size.
func CalculatePaginationMetadata(totalRecords, page, pageSize int) PaginationMetadata {
	if totalRecords == 0 {
		// Note that we return an empty Metadata struct if there are no records.
		return PaginationMetadata{}
	}
	return PaginationMetadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

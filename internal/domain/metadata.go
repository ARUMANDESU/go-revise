package domain

import "math"

// PaginationMetadata represents the metadata for paginated responses.
type PaginationMetadata struct {
	CurrentPage  int32 `json:"current_page"`
	PageSize     int32 `json:"page_size"`
	FirstPage    int32 `json:"first_page"`
	LastPage     int32 `json:"last_page"`
	TotalRecords int32 `json:"total_records"`
}

// CalculatePaginationMetadata calculates the pagination metadata based on the total number of records, the current page, and the page size.
func CalculatePaginationMetadata(totalRecords, page, pageSize int32) PaginationMetadata {
	if totalRecords == 0 {
		// Note that we return an empty Metadata struct if there are no records.
		return PaginationMetadata{}
	}
	return PaginationMetadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int32(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

package domain

const (
	PaginationDefaultPage     = 1
	PaginationDefaultPageSize = 10
)

// Pagination is the pagination configuration.
// Recommended to use `NewPagination` to create a new `Pagination` struct, because it sets default values.
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func (p Pagination) Limit() int {
	return p.PageSize
}
func (p Pagination) Offset() int {
	// handle negative values because client might not use `NewPagination`
	if p.Page <= 0 || p.PageSize <= 0 {
		return 0
	}
	return (p.Page - 1) * p.PageSize
}

// NewPagination creates a new `Pagination` struct with default values if none are provided.
func NewPagination(page, pageSize int) *Pagination {
	if page <= 0 {
		page = PaginationDefaultPage
	}
	if pageSize <= 0 {
		pageSize = PaginationDefaultPageSize
	}

	return &Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func DefaultPagination() *Pagination {
	return NewPagination(PaginationDefaultPage, PaginationDefaultPageSize)
}

// SortField is the field to sort by.
// Be sure that this fields are the same as the fields in the database.
type SortField string

const (
	// This is used if the client doesn't provide the sort field.
	SortFieldDefault        SortField = "next_revision_at"
	SortFieldID             SortField = "id"
	SortFieldCreatedAt      SortField = "created_at"
	SortFieldUpdatedAt      SortField = "updated_at"
	SortFieldTags           SortField = "tags"
	SortFieldLastRevisedAt  SortField = "last_revised_at"
	SortFieldNextRevisionAt SortField = "next_revision_at"
)

type SortOrder string

const (
	// This is used if the client doesn't provide the sort order.
	SortOrderDefault SortOrder = "desc"
	SortOrderAsc     SortOrder = "asc"
	SortOrderDesc    SortOrder = "desc"
)

// Sort is the sorting configuration.
// Recommended to use NewSort to create a new Sort struct, because it sets default values.
type Sort struct {
	// Field is the field to sort by.
	Field SortField `json:"field"`
	Order SortOrder `json:"order"`
}

// NewSort creates a new `Sort` struct with default values if none are provided.
func NewSort(field SortField, order SortOrder) *Sort {
	if field == "" {
		field = SortFieldDefault
	}
	if order == "" {
		order = SortOrderDefault
	}
	return &Sort{
		Field: field,
		Order: order,
	}
}

func DefaultSort() *Sort {
	return NewSort(SortFieldDefault, SortOrderDefault)
}

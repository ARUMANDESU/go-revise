package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagination_Limit(t *testing.T) {
	tests := []struct {
		name     string
		pageSize int
		want     int
	}{
		{
			name:     "default page size",
			pageSize: PaginationDefaultPageSize,
			want:     PaginationDefaultPageSize,
		},
		{
			name:     "custom page size",
			pageSize: 20,
			want:     20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{PageSize: tt.pageSize}
			assert.Equal(t, tt.want, p.Limit())
		})
	}
}

func TestPagination_Offset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		want     int
	}{
		{
			name:     "first page",
			page:     1,
			pageSize: 10,
			want:     0,
		},
		{
			name:     "second page",
			page:     2,
			pageSize: 10,
			want:     10,
		},
		{
			name:     "third page",
			page:     3,
			pageSize: 5,
			want:     10,
		},
		{
			name:     "0 page",
			page:     1,
			pageSize: 10,
			want:     0,
		},
		{
			name:     "negative page and page size",
			page:     -1,
			pageSize: -1,
			want:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{Page: tt.page, PageSize: tt.pageSize}
			assert.Equal(t, tt.want, p.Offset())
		})
	}
}

func TestNewPagination(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		want     Pagination
	}{
		{
			name:     "default values",
			page:     0,
			pageSize: 0,
			want:     Pagination{Page: PaginationDefaultPage, PageSize: PaginationDefaultPageSize},
		},
		{
			name:     "custom values",
			page:     2,
			pageSize: 20,
			want:     Pagination{Page: 2, PageSize: 20},
		},
		{
			name:     "negative values",
			page:     -1,
			pageSize: -1,
			want:     Pagination{Page: PaginationDefaultPage, PageSize: PaginationDefaultPageSize},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pagination := NewPagination(tt.page, tt.pageSize)
			assert.Equal(t, tt.want, *pagination)
		})
	}
}

func TestNewSort(t *testing.T) {
	tests := []struct {
		name  string
		field SortField
		order SortOrder
		want  Sort
	}{
		{
			name:  "default values",
			field: "",
			order: "",
			want:  Sort{Field: SortFieldDefault, Order: SortOrderDefault},
		},
		{
			name:  "custom values",
			field: SortFieldCreatedAt,
			order: SortOrderAsc,
			want:  Sort{Field: SortFieldCreatedAt, Order: SortOrderAsc},
		},
		{
			name:  "default field, custom order",
			field: "",
			order: SortOrderAsc,
			want:  Sort{Field: SortFieldDefault, Order: SortOrderAsc},
		},
		{
			name:  "custom field, default order",
			field: SortFieldUpdatedAt,
			order: "",
			want:  Sort{Field: SortFieldUpdatedAt, Order: SortOrderDefault},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort := NewSort(tt.field, tt.order)
			assert.Equal(t, tt.want, *sort)
		})
	}
}

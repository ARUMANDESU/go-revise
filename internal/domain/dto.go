package domain

type CreateReviseItemDTO struct {
	UserID      string   `json:"user_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type UpdateReviseItemDTO struct {
	ID          string   `json:"id"`
	UserID      string   `json:"user_id"` // who is updating the revise
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	// Fields that must be updated should be included in this list
	// needed for partial update
	UpdateFields []string `json:"update_fields"`
}

type ListReviseItemDTO struct {
	// UserID is the user who is listing the items.
	// Allowed types: string(uuid) and int64 (telegram user id).
	// Type is any because it can be either string or int64.
	UserID any `json:"user_id"`
	// Pagination is the pagination configuration.
	Pagination *Pagination `json:"pagination"` // not required, that's why it's a pointer
	// Sort is the sorting configuration.
	Sort *Sort `json:"sort"` // not required, that's why it's a pointer
}

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

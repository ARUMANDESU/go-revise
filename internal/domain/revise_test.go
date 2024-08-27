package domain

import (
	"database/sql/driver"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringArray_Scan(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    StringArray
		wantErr bool
	}{
		{
			name:    "nil value",
			input:   nil,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "valid string",
			input:   "tag1,tag2,tag3",
			want:    StringArray{"tag1", "tag2", "tag3"},
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "invalid type",
			input:   123,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a StringArray
			err := a.Scan(tt.input)
			require.Equal(t, tt.wantErr, err != nil)

			assert.Equal(t, tt.want, a)
		})
	}
}

func TestStringArray_Value(t *testing.T) {
	tests := []struct {
		name string
		a    StringArray
		want driver.Value
	}{
		{
			name: "multiple tags",
			a:    StringArray{"tag1", "tag2", "tag3"},
			want: "tag1,tag2,tag3",
		},
		{
			name: "single tag",
			a:    StringArray{"tag1"},
			want: "tag1",
		},
		{
			name: "empty array",
			a:    StringArray{},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.a.Value())
		})
	}
}

func TestReviseItem_AbleToUpdate(t *testing.T) {
	uid, _ := uuid.NewV7()
	anotherUID, _ := uuid.NewV7()

	tests := []struct {
		name string
		r    ReviseItem
		id   string
		want bool
	}{
		{
			name: "able to update",
			r:    ReviseItem{UserID: uid},
			id:   uid.String(),
			want: true,
		},
		{
			name: "unable to update",
			r:    ReviseItem{UserID: uid},
			id:   anotherUID.String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.r.AbleToUpdate(tt.id))
		})
	}
}

func TestReviseItem_PartialUpdate(t *testing.T) {
	revise := ReviseItem{
		Name:        "old name",
		Description: "old description",
		Tags:        StringArray{"old tag1", "old tag2"},
	}

	tests := []struct {
		name string
		dto  UpdateReviseItemDTO
		want ReviseItem
	}{
		{
			name: "update name",
			dto: UpdateReviseItemDTO{
				Name:         " new name ",
				Description:  "old description ",
				Tags:         StringArray{"old tag1", "old tag2"},
				UpdateFields: []string{"name"},
			},
			want: ReviseItem{
				Name:        "new name",
				Description: "old description",
				Tags:        StringArray{"old tag1", "old tag2"},
			},
		},
		{
			name: "update description",
			dto: UpdateReviseItemDTO{
				Name:         "old name",
				Description:  "new description",
				Tags:         StringArray{"old tag1", "old tag2"},
				UpdateFields: []string{"description"},
			},
			want: ReviseItem{
				Name:        "old name",
				Description: "new description",
				Tags:        StringArray{"old tag1", "old tag2"},
			},
		},
		{
			name: "update tags",
			dto: UpdateReviseItemDTO{
				Name:         "old name",
				Description:  "old description",
				Tags:         StringArray{"new tag1", "new tag2"},
				UpdateFields: []string{"tags"},
			},
			want: ReviseItem{
				Name:        "old name",
				Description: "old description",
				Tags:        StringArray{"new tag1", "new tag2"},
			},
		},
		{
			name: "update all",
			dto: UpdateReviseItemDTO{
				Name:         "new name",
				Description:  "new description",
				Tags:         StringArray{"new tag1", "new tag2"},
				UpdateFields: []string{"name", "description", "tags"},
			},
			want: ReviseItem{
				Name:        "new name",
				Description: "new description",
				Tags:        StringArray{"new tag1", "new tag2"},
			},
		},
		{
			name: "no update",
			dto: UpdateReviseItemDTO{
				Name:         "old name",
				Description:  "old description",
				Tags:         StringArray{"old tag1", "old tag2"},
				UpdateFields: []string{},
			},
			want: ReviseItem{
				Name:        "old name",
				Description: "old description",
				Tags:        StringArray{"old tag1", "old tag2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated := revise.PartialUpdate(tt.dto)

			for _, field := range tt.dto.UpdateFields {
				switch field {
				case "name":
					assert.Equal(t, tt.want.Name, updated.Name, "name field")
				case "description":
					assert.Equal(t, tt.want.Description, updated.Description, "description field")
				case "tags":
					assert.Equal(t, tt.want.Tags, updated.Tags, "tags field")
				}
			}

			if tt.dto.UpdateFields != nil || len(tt.dto.UpdateFields) > 0 {
				assert.NotEqual(t, revise.UpdatedAt, updated.UpdatedAt)
			} else {
				assert.Equal(t, revise.UpdatedAt, updated.UpdatedAt)
			}
		})
	}
}

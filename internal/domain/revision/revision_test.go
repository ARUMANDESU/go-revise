package revision

import (
	"testing"

	"github.com/gofrs/uuid"
)

func TestNewRevision(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		reviseItemID uuid.UUID
		wantErr      bool
	}{
		{
			name:         "With valid revise item ID",
			reviseItemID: uuid.Must(uuid.NewV7()),
			wantErr:      false,
		},
		{
			name:         "With invalid revise item ID",
			reviseItemID: uuid.Nil,
			wantErr:      true,
		},
		{
			name:         "With empty revise item ID",
			reviseItemID: uuid.UUID{},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewRevision(tt.reviseItemID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRevision() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

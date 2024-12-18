package reviseitem

import (
	"strings"
	"testing"
	"time"

	"github.com/clarify/subtest"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	domainUser "github.com/ARUMANDESU/go-revise/internal/domain/user"
	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
)

func TestNewReviseItem(t *testing.T) {
	t.Parallel()

	reviseItemID := NewReviseItemID()
	userID := domainUser.NewUserID()

	tests := []struct {
		name    string
		args    NewReviseItemArgs
		want    *ReviseItem
		wantErr bool
	}{
		{
			name: "With valid arguments",
			args: validNewReviseItemArgs(t, reviseItemID, userID),
			want: &ReviseItem{
				id:             reviseItemID,
				userID:         userID,
				name:           validName(t, language.Kazakh),
				description:    validDescription(t, language.Kazakh),
				tags:           validTags(t),
				createdAt:      time.Now(),
				updatedAt:      time.Now(),
				deletedAt:      nil,
				nextRevisionAt: validNextRevisionAt(t),
				lastRevisedAt:  time.Time{},
			},
		},
		{
			name:    "With invalid arguments",
			args:    invalidNewReviseItemArgs(t, reviseItemID, userID),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewReviseItem(tt.args)
			if tt.wantErr {
				t.Run("Expect error", subtest.Value(err).Error())
			} else {
				t.Run("Expect no error", subtest.Value(err).NoError())
				t.Run("Expect revise item be equal", func(t *testing.T) {
					assertReviseItem(t, got, tt.want)
				})
			}
		})
	}
}

func TestReviseItem_UpdateName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		newName string
		wantErr bool
	}{
		{
			name:    "With valid name",
			newName: validName(t, language.Kazakh),
			wantErr: false,
		},
		{
			name:    "With long name",
			newName: longName(t),
			wantErr: true,
		},
		{
			name:    "With empty name",
			newName: "",
			wantErr: true,
		},
		{
			name:    "With name in Russian",
			newName: validName(t, language.Russian),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reviseItem := validReviseItem(t)

			err := reviseItem.UpdateName(tt.newName)
			if tt.wantErr {
				t.Run("Expect error", subtest.Value(err).Error())
			} else {
				t.Run("Expect no error", subtest.Value(err).NoError())
				t.Run("Expect name to be updated", func(t *testing.T) {
					tt.newName = strings.TrimSpace(tt.newName)
					assert.Equal(t, tt.newName, reviseItem.name)
				})
			}
		})
	}
}

func TestReviseItem_UpdateDescription(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		newDescription string
		wantErr        bool
	}{
		{
			name:           "With valid description",
			newDescription: validDescription(t, language.Kazakh),
			wantErr:        false,
		},
		{
			name:           "With long description",
			newDescription: longDescription(t),
			wantErr:        true,
		},
		{
			name:           "With empty description",
			newDescription: "",
			wantErr:        false,
		},
		{
			name:           "With description in Russian",
			newDescription: validDescription(t, language.Russian),
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reviseItem := validReviseItem(t)

			err := reviseItem.UpdateDescription(tt.newDescription)
			if tt.wantErr {
				t.Run("Expect error", subtest.Value(err).Error())
			} else {
				t.Run("Expect no error", subtest.Value(err).NoError())
				t.Run("Expect description to be updated", func(t *testing.T) {
					tt.newDescription = strings.TrimSpace(tt.newDescription)
					assert.Equal(t, tt.newDescription, reviseItem.description)
				})
			}
		})
	}
}

func TestReviseItem_AddTags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		newTags valueobject.Tags
		wantErr bool
	}{
		{
			name:    "With valid tags",
			newTags: validTags(t),
			wantErr: false,
		},
		{
			name:    "With empty tags",
			newTags: valueobject.NewTags(),
			wantErr: true,
		},
		{
			name:    "With a lot of tags",
			newTags: valueobject.ExampleInValidTags(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reviseItem := validReviseItem(t)

			err := reviseItem.AddTags(tt.newTags)
			if tt.wantErr {
				t.Run("Expect error", subtest.Value(err).Error())
			} else {
				t.Run("Expect no error", subtest.Value(err).NoError())
				t.Run("Expect tags to be updated", func(t *testing.T) {
					tt.newTags = tt.newTags.TrimSpace()
					assert.Equal(t, tt.newTags, reviseItem.tags)
				})
			}
		})
	}
}

func TestReviseItem_RemoveTags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		reviseItem *ReviseItem
		tags       valueobject.Tags
		expected   valueobject.Tags
		wantErr    bool
	}{
		{
			name:       "With valid tags",
			reviseItem: validReviseItem(t),
			tags:       validTags(t),
			expected:   valueobject.NewTags(),
			wantErr:    false,
		},
		{
			name:       "With empty tags",
			reviseItem: validReviseItem(t),
			tags:       valueobject.NewTags(),
			expected:   valueobject.NewTags(),
			wantErr:    true,
		},
		{
			name: "With invalid long tag",
			reviseItem: &ReviseItem{
				tags: valueobject.NewTags(valueobject.ExampleInvalidTag),
			},
			tags:     valueobject.NewTags(valueobject.ExampleInvalidTag),
			expected: valueobject.NewTags(),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.reviseItem.RemoveTags(tt.tags)
			if tt.wantErr {
				t.Run("Expect error", subtest.Value(err).Error())
			} else {
				t.Run("Expect no error", subtest.Value(err).NoError())
				t.Run("Expect tags to be updated", func(t *testing.T) {
					assertTags(t, tt.reviseItem.tags, tt.expected)
				})
			}
		})
	}
}

func TestReviseItem_UpdateNextRevisionAt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		nextRevisionAt time.Time
		wantErr        bool
	}{
		{
			name:           "With valid next revision at",
			nextRevisionAt: validNextRevisionAt(t),
			wantErr:        false,
		},
		{
			name:           "With invalid next revision at",
			nextRevisionAt: time.Time{},
			wantErr:        true,
		},
		{
			name:           "With past next revision at",
			nextRevisionAt: time.Now().Add(-time.Hour),
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reviseItem := validReviseItem(t)

			err := reviseItem.UpdateNextRevisionAt(tt.nextRevisionAt)
			if tt.wantErr {
				t.Run("Expect error", subtest.Value(err).Error())
			} else {
				t.Run("Expect no error", subtest.Value(err).NoError())
				t.Run("Expect next revision at to be updated", func(t *testing.T) {
					assert.WithinDuration(t, tt.nextRevisionAt, reviseItem.nextRevisionAt, time.Second)
				})
			}
		})
	}
}

func TestReviseItem_MarkAsDeleted(t *testing.T) {
	t.Parallel()

	reviseItem := validReviseItem(t)

	reviseItem.MarkAsDeleted()

	t.Run("Expect deletedAt to be set", func(t *testing.T) {
		assert.NotNil(t, reviseItem.deletedAt)
		assert.WithinDuration(t, time.Now(), *reviseItem.deletedAt, time.Second)
		assert.WithinDuration(t, time.Now(), reviseItem.updatedAt, time.Second)
	})
}

func TestReviseItem_Restore(t *testing.T) {
	t.Parallel()

	reviseItem := validReviseItem(t)

	reviseItem.MarkAsDeleted()
	reviseItem.Restore()

	t.Run("Expect deletedAt to be nil", func(t *testing.T) {
		assert.Nil(t, reviseItem.deletedAt)
		assert.WithinDuration(t, time.Now(), reviseItem.updatedAt, time.Second)
	})
}

func TestReviseItem_CanModify(t *testing.T) {
	t.Parallel()

	userID := domainUser.NewUserID()
	otherUserID := domainUser.NewUserID()

	tests := []struct {
		name       string
		reviseItem *ReviseItem
		userID     uuid.UUID
		want       bool
	}{
		{
			name: "With matching userID",
			reviseItem: &ReviseItem{
				userID: userID,
			},
			userID: userID,
			want:   true,
		},
		{
			name: "With non-matching userID",
			reviseItem: &ReviseItem{
				userID: userID,
			},
			userID: otherUserID,
			want:   false,
		},
		{
			name: "With nil userID",
			reviseItem: &ReviseItem{
				userID: userID,
			},
			userID: uuid.UUID{},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.reviseItem.CanModify(tt.userID)
			assert.Equal(t, tt.want, got)
		})
	}
}

func assertReviseItem(t *testing.T, got, want *ReviseItem) {
	t.Helper()

	assert.Equal(t, got.id, want.id)
	assert.Equal(t, got.userID, want.userID)
	assert.Equal(t, got.name, want.name)
	assert.Equal(t, got.description, want.description)
	assertTags(t, got.tags, want.tags)
	assert.WithinDuration(t, got.createdAt, want.createdAt, time.Second)
	assert.WithinDuration(t, got.updatedAt, want.updatedAt, time.Second)
	assert.Equal(t, got.deletedAt, want.deletedAt)
	assert.WithinDuration(t, got.nextRevisionAt, want.nextRevisionAt, time.Second)
	assert.WithinDuration(t, got.lastRevisedAt, want.lastRevisedAt, time.Second)
}

func assertTags(t *testing.T, got, want valueobject.Tags) {
	t.Helper()

	if !valueobject.IsTagsEqual(&got, &want) {
		t.Errorf("got tags: %v, want tags: %v", got, want)
	}
}

func validReviseItem(t *testing.T) *ReviseItem {
	t.Helper()

	reviseItemID := NewReviseItemID()
	userID := domainUser.NewUserID()

	reviseItem, err := NewReviseItem(validNewReviseItemArgs(t, reviseItemID, userID))
	assert.NoError(t, err)

	return reviseItem
}

func validNewReviseItemArgs(t *testing.T, reviseID, userID uuid.UUID) NewReviseItemArgs {
	t.Helper()
	return NewReviseItemArgs{
		ID:          reviseID,
		UserID:      userID,
		Name:        validName(t, language.Kazakh),
		Description: validDescription(t, language.Kazakh),
		Tags:        validTags(t),
	}
}

func invalidNewReviseItemArgs(t *testing.T, reviseID, userID uuid.UUID) NewReviseItemArgs {
	t.Helper()
	return NewReviseItemArgs{
		ID:          reviseID,
		UserID:      userID,
		Name:        longName(t),
		Description: longDescription(t),
		Tags:        valueobject.ExampleInValidTags(),
	}
}

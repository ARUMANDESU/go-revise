package reviseitem

import (
	"github.com/ARUMANDESU/go-revise/internal/domain/revision"
)

// Aggregate represents a revise item aggregate.
type Aggregate struct {
	ReviseItem
	revisions []revision.Revision
}

func NewAggregate(item *ReviseItem) *Aggregate {
	return &Aggregate{ReviseItem: *item}
}

func (a *Aggregate) Review() {
	a.revisions = append(a.revisions, *revision.NewRevision())
}

func (a *Aggregate) Revisions() []revision.Revision {
	return a.revisions
}

package reviseitem

import (
	"github.com/ARUMANDESU/go-revise/internal/domain/revision"
)

// Aggregate represents a revise item aggregate.
type Aggregate struct {
	ReviseItem
	Revisions []revision.Revision
}

func (a *Aggregate) AddRevision(rev revision.Revision) {
	a.Revisions = append(a.Revisions, rev)
}

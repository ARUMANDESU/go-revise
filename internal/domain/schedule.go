package domain

import "time"

type ScheduledItem struct {
	User
	ReviseItem
}

func (s ScheduledItem) NotifyAt() time.Time {
	if s.NextRevisionAt.Before(time.Now()) {
		return time.Now().Add(1 * time.Second)
	}

	return s.NextRevisionAt
}

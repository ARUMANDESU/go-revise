package valueobject

import (
	"strconv"
	"strings"
	"time"

	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

const maxReviewIntervals = 17

const (
	minuteCharacter = "m"
	hourCharacter   = "h"
	dayCharacter    = "d"
	weekCharacter   = "w"
	monthCharacter  = "M"
	yearCharacter   = "y"
)

// ReviewInterval represents the interval of time between reviews.
type ReviewInterval struct {
	value [maxReviewIntervals]time.Duration
}

func (r ReviewInterval) Next(i int) time.Time {
	var add time.Duration
	if i < 0 || i >= maxReviewIntervals {
		return time.Time{}
	}
	add = r.value[i]

	return time.Now().Add(add)
}

// Default review intervals.
var defaultIntervals = [maxReviewIntervals]time.Duration{
	time.Hour * 24,           // 1 day
	time.Hour * 24 * 3,       // 3 days
	time.Hour * 24 * 7,       // 1 week
	time.Hour * 24 * 7 * 2,   // 2 weeks
	time.Hour * 24 * 7 * 3,   // 3 weeks
	time.Hour * 24 * 30,      // 1 month
	time.Hour * 24 * 45,      // 1 month and a half
	time.Hour * 24 * 60,      // 2 months
	time.Hour * 24 * 90,      // 3 months
	time.Hour * 24 * 120,     // 4 months
	time.Hour * 24 * 180,     // 6 months
	time.Hour * 24 * 270,     // 9 months
	time.Hour * 24 * 365,     // 1 year
	time.Hour * 24 * 30 * 18, // 1 year and a half
	time.Hour * 24 * 365 * 2, // 2 years
	time.Hour * 24 * 365 * 3, // 3 years
	time.Hour * 24 * 365 * 5, // 5 years
}

func DefaultReviewIntervals() ReviewInterval {
	return ReviewInterval{value: defaultIntervals}
}

// ParseReviewInterval parses the review interval from a string separated by spaces.
func ParseReviewInterval(interval string) (ReviewInterval, error) {
	op := errs.Op("valueobject.parse_review_interval")
	intervals := [maxReviewIntervals]time.Duration{}
	split := splitInterval(interval)

	if len(split) > maxReviewIntervals {
		return ReviewInterval{}, errs.
			NewIncorrectInputError(op, errs.ErrInvalidInput, "too many intervals").
			WithMessages([]errs.Message{{Key: "message", Value: "too many intervals"}}).
			WithContext("intervals", interval)
	}

	for i, v := range split {
		duration, err := parseDuration(v)
		if err != nil {
			return ReviewInterval{}, errs.WithOp(op, err, "failed to parse duration")
		}
		intervals[i] = duration
	}

	if len(split) < maxReviewIntervals {
		copy(intervals[len(split):], defaultIntervals[len(split):])
	}

	return ReviewInterval{
		value: intervals,
	}, nil
}

func splitInterval(interval string) []string {
	return strings.Split(interval, " ")
}

// parseDuration parses the duration from a string with a number followed by a time unit.
func parseDuration(interval string) (time.Duration, error) {
	op := errs.Op("valueobject.parse_duration")
	value := interval[:len(interval)-1]
	unit := interval[len(interval)-1:]

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, errs.
			NewIncorrectInputError(op, err, "invalid value").
			WithMessages([]errs.Message{{Key: "message", Value: "invalid value"}}).
			WithContext("value", value)
	}

	var duration time.Duration
	switch unit {
	case minuteCharacter:
		duration = time.Duration(floatValue * float64(time.Minute))
	case hourCharacter:
		duration = time.Duration(floatValue * float64(time.Hour))
	case dayCharacter:
		duration = time.Duration(floatValue * float64(time.Hour) * 24)
	case weekCharacter:
		duration = time.Duration(floatValue * float64(time.Hour) * 24 * 7)
	case monthCharacter:
		duration = time.Duration(floatValue * float64(time.Hour) * 24 * 30)
	case yearCharacter:
		duration = time.Duration(floatValue * float64(time.Hour) * 24 * 365)
	default:
		return 0, errs.
			NewIncorrectInputError(op, errs.ErrInvalidInput, "invalid unit").
			WithMessages([]errs.Message{{Key: "message", Value: "invalid unit"}}).
			WithContext("unit", unit)
	}

	return duration, nil
}

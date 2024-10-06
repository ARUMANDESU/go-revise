package valueobject

import (
	"fmt"
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
	intervals := [maxReviewIntervals]time.Duration{}
	split := splitInterval(interval)

	if len(split) > maxReviewIntervals {
		return ReviewInterval{}, errs.NewIncorrectInputError(fmt.Sprintf("too many values provided, max is %d", maxReviewIntervals), "too-many-values")
	}

	for i, v := range split {
		duration, err := parseDuration(v)
		if err != nil {
			return ReviewInterval{}, err
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
	value := interval[:len(interval)-1]
	unit := interval[len(interval)-1:]

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, errs.NewIncorrectInputError(fmt.Sprintf("invalid value type: %s", value), "invalid-value-type")
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
		return 0, errs.NewIncorrectInputError(fmt.Sprintf("invalid time unit: %s", unit), "invalid-time-unit")
	}

	return duration, nil
}

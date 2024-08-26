package domain

import "time"

// ReviseInterval is a type that represents the interval at which a ReviseItem should be revised.
type ReviseInterval time.Duration

// ReviseIteration is a type that represents the number of iterations that a ReviseItem has gone through.
type ReviseIteration uint

const (
	OneDay         ReviseInterval = ReviseInterval(24 * time.Hour)
	ThreeDays      ReviseInterval = ReviseInterval(3 * 24 * time.Hour)
	OneWeek        ReviseInterval = ReviseInterval(7 * 24 * time.Hour)
	TwoWeeks       ReviseInterval = ReviseInterval(14 * 24 * time.Hour)
	ThreeWeeks     ReviseInterval = ReviseInterval(21 * 24 * time.Hour)
	OneMonth       ReviseInterval = ReviseInterval(30 * 24 * time.Hour)
	OneHalfMonth   ReviseInterval = ReviseInterval(45 * 24 * time.Hour)
	TwoMonths      ReviseInterval = ReviseInterval(60 * 24 * time.Hour)
	ThreeMonths    ReviseInterval = ReviseInterval(90 * 24 * time.Hour)
	FourMonths     ReviseInterval = ReviseInterval(120 * 24 * time.Hour)
	SixMonths      ReviseInterval = ReviseInterval(180 * 24 * time.Hour)
	NineMonths     ReviseInterval = ReviseInterval(270 * 24 * time.Hour)
	OneYear        ReviseInterval = ReviseInterval(365 * 24 * time.Hour)
	EighteenMonths ReviseInterval = ReviseInterval(18 * 30 * 24 * time.Hour)
	TwoYears       ReviseInterval = ReviseInterval(2 * 365 * 24 * time.Hour)
	ThreeYears     ReviseInterval = ReviseInterval(3 * 365 * 24 * time.Hour)
	FiveYears      ReviseInterval = ReviseInterval(5 * 365 * 24 * time.Hour)
)

var IntervalMap = map[ReviseIteration]ReviseInterval{
	0:  OneDay,
	1:  OneDay,
	2:  ThreeDays,
	3:  OneWeek,
	4:  TwoWeeks,
	5:  ThreeWeeks,
	6:  OneMonth,
	7:  OneHalfMonth,
	8:  TwoMonths,
	9:  ThreeMonths,
	10: FourMonths,
	11: SixMonths,
	12: NineMonths,
	13: OneYear,
	14: EighteenMonths,
	15: TwoYears,
	16: ThreeYears,
	17: FiveYears,
}

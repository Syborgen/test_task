package datastructures

import (
	"fmt"
	"strings"
	"time"
)

type TimeRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

const timeFormat = "2006-01-02 15:04:05"

// time range as string pattern: ["YYYY-MM-DD HH:MI:SS","YYYY-MM-DD HH:MI:SS")
func NewTimeRange(timeRangeAsString string) TimeRange {
	trimmedString := strings.Trim(timeRangeAsString, "[\")")
	splittedString := strings.Split(trimmedString, "\",\"")
	return TimeRange{Start: splittedString[0], End: splittedString[1]}
}

func (thisTR TimeRange) IsContains(tr TimeRange) bool {
	thisStart, err := time.Parse(timeFormat, thisTR.Start)
	if err != nil {
		return false
	}

	thisEnd, err := time.Parse(timeFormat, thisTR.End)
	if err != nil {
		return false
	}

	dStart, err := time.Parse(timeFormat, tr.Start)
	if err != nil {
		return false
	}

	dEnd, err := time.Parse(timeFormat, tr.End)
	if err != nil {
		return false
	}

	return (thisStart.Before(dStart) || thisStart.Equal(dStart)) && thisEnd.After(dEnd)
}

const dateAsStringLength = 19

func (thisTR TimeRange) IsOverlapped(tr TimeRange) bool {
	thisStart, err := time.Parse(timeFormat, thisTR.Start[:dateAsStringLength])
	if err != nil {
		return false
	}

	thisEnd, err := time.Parse(timeFormat, thisTR.End[:dateAsStringLength])
	if err != nil {
		return false
	}

	trStart, err := time.Parse(timeFormat, tr.Start[:dateAsStringLength])
	if err != nil {
		return false
	}

	trEnd, err := time.Parse(timeFormat, tr.End[:dateAsStringLength])
	if err != nil {
		return false
	}

	return !(trEnd.Before(thisStart) || trEnd.Equal(thisStart)) &&
		!(trStart.After(thisEnd) || trStart.Equal(thisEnd))
}

func (tr TimeRange) GetDuration() time.Duration {
	start, err := time.Parse(timeFormat, tr.Start)
	if err != nil {
		return time.Duration(0)
	}

	end, err := time.Parse(timeFormat, tr.End)
	if err != nil {
		return time.Duration(0)
	}

	return end.Sub(start)
}

func (tr TimeRange) String() string {
	return fmt.Sprintf("[%s, %s)", tr.Start, tr.End)
}

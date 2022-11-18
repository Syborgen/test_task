package datastructures

import (
	"errors"
	"fmt"
	"time"
)

const serverLocationName = "Europe/Moscow"

const TimeFormat = "2006-01-02 15:04:05"

type TimeRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func (tr *TimeRange) String() string {
	return fmt.Sprintf("[%s, %s)", tr.Start, tr.End)
}

func NewTimeRange(start string, end string) (TimeRange, error) {
	startTime, err := time.Parse(TimeFormat, start)
	if err != nil {
		return TimeRange{}, fmt.Errorf("parse error: %w", err)
	}

	endTime, err := time.Parse(TimeFormat, end)
	if err != nil {
		return TimeRange{}, fmt.Errorf("parse error: %w", err)
	}

	if !startTime.Before(endTime) {
		return TimeRange{}, errors.New("start must be before end")
	}

	return TimeRange{Start: startTime.Format(TimeFormat), End: endTime.Format(TimeFormat)}, nil
}

func (tr *TimeRange) ConvertToLocalTime() error {
	serverLocation, err := time.LoadLocation(serverLocationName)
	if err != nil {
		return fmt.Errorf("load location error: %w", err)
	}

	tr.Start, err = convertTimeFromLocationToLocation(tr.Start, serverLocation, time.Now().Location())
	if err != nil {
		return fmt.Errorf("convert 'start' to server time error: %w", err)
	}

	tr.End, err = convertTimeFromLocationToLocation(tr.End, serverLocation, time.Now().Location())
	if err != nil {
		return fmt.Errorf("convert 'end' to server time error: %w", err)
	}

	return nil
}

func (tr *TimeRange) ConvertToServerTime() error {
	serverLocation, err := time.LoadLocation(serverLocationName)
	if err != nil {
		return fmt.Errorf("load location error: %w", err)
	}

	tr.Start, err = convertTimeFromLocationToLocation(tr.Start, time.Now().Location(), serverLocation)
	if err != nil {
		return fmt.Errorf("convert 'start' to server time error: %w", err)
	}

	tr.End, err = convertTimeFromLocationToLocation(tr.End, time.Now().Location(), serverLocation)
	if err != nil {
		return fmt.Errorf("convert 'end' to server time error: %w", err)
	}

	return nil
}

func convertTimeFromLocationToLocation(timeAsString string, from *time.Location, to *time.Location) (string, error) {
	timeToConvert, err := time.ParseInLocation(TimeFormat, timeAsString, from)
	if err != nil {
		return "", fmt.Errorf("parse error: %w", err)
	}

	convertedTime := timeToConvert.In(to)
	return convertedTime.Format(TimeFormat), nil
}

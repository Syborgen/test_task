package datastructures

import (
	"testing"
	"time"
)

func TestConvertToServerTime(t *testing.T) {
	expected := "2022-01-01 15:00:00"

	input := "2022-01-01 12:00:00"
	clientLocation := time.UTC

	serverLocation, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		t.Errorf("load location error: %s", err)
		return
	}

	actual, err := convertTimeFromLocationToLocation(input, clientLocation, serverLocation)
	if err != nil {
		t.Errorf("convert error: %s", err)
		return
	}

	if expected != actual {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

func TestConvertToLocalTime(t *testing.T) {
	expected := "2022-01-01 12:00:00"

	input := "2022-01-01 15:00:00"
	clientLocation := time.Local

	serverLocation, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		t.Errorf("load location error: %s", err)
		return
	}

	actual, err := convertTimeFromLocationToLocation(input, serverLocation, clientLocation)
	if err != nil {
		t.Errorf("convert error: %s", err)
		return
	}

	if expected != actual {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

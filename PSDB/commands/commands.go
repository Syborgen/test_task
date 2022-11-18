package commands

import (
	"PSBD/datastructures"
	"PSBD/dbhelper"
	"errors"
	"fmt"
	"sort"
	"time"
)

func GetObjects() ([]datastructures.Object, error) {
	return dbhelper.SelectObject(dbhelper.DB)
}

func GetTechWindows() ([]datastructures.GroupedTechWindow, error) {
	return dbhelper.SelectGroupedTechWindows(dbhelper.DB)
}

func GetTechWindowsAll() ([]datastructures.TechWindow, error) {
	return dbhelper.SelectTechWindows(dbhelper.DB)
}

func GetSortedWindowsQuery(sort string, start string, end string) ([]datastructures.GroupedTechWindow, error) {
	return dbhelper.SelectSortedTechWindowWithBounds(dbhelper.DB, start, end, sort)
}

func GetSortedWindowsProc(sort string, start string, end string) ([]datastructures.GroupedTechWindow, error) {
	techWindows, err := dbhelper.SelectTechWindows(dbhelper.DB)
	if err != nil {
		return nil, fmt.Errorf("select tech window error: %w", err)
	}
	bounds := datastructures.TimeRange{Start: start, End: end}
	filteredTechWindows := filterTechWindowsByTimeBounds(techWindows, bounds)
	groupedTechWindows := groupTechWindows(filteredTechWindows)
	sortedTechWindows := sortGroupedTechWindows(groupedTechWindows, sort)
	return sortedTechWindows, nil
}

func filterTechWindowsByTimeBounds(techWindows []datastructures.TechWindow,
	timeBound datastructures.TimeRange) []datastructures.TechWindow {
	result := []datastructures.TechWindow{}
	for _, techWindow := range techWindows {
		if timeBound.IsContains(techWindow.Duration) {
			result = append(result, techWindow)
		}
	}

	return result
}

func groupTechWindows(techWindows []datastructures.TechWindow) []datastructures.GroupedTechWindow {
	groupsByID := make(map[int][]datastructures.TechWindow)
	for _, techWindow := range techWindows {
		groupsByID[techWindow.IDObject] = append(groupsByID[techWindow.IDObject], techWindow)
	}

	result := []datastructures.GroupedTechWindow{}
	for idObject, groupByID := range groupsByID {
		count := len(groupByID)
		averageDuration := calculateAverageDuration(groupByID)
		result = append(result, datastructures.GroupedTechWindow{
			IDObject:        idObject,
			WindowsCount:    count,
			AverageDuration: averageDuration.String(),
		})
	}

	return result
}

func calculateAverageDuration(groupByID []datastructures.TechWindow) time.Duration {
	sumOfDurations := time.Duration(0)
	for _, techWindow := range groupByID {
		sumOfDurations += techWindow.Duration.GetDuration()
	}

	averageDuration := sumOfDurations / time.Duration(len(groupByID))
	return averageDuration
}

func sortGroupedTechWindows(groupedTechWindows []datastructures.GroupedTechWindow,
	sortOrder string) []datastructures.GroupedTechWindow {
	switch sortOrder {
	case "asc":
		sort.Slice(groupedTechWindows, func(i int, j int) bool {
			durationForI, err := time.ParseDuration(groupedTechWindows[i].AverageDuration)
			if err != nil {
				return false
			}

			durationForJ, err := time.ParseDuration(groupedTechWindows[j].AverageDuration)
			if err != nil {
				return false
			}

			return durationForI < durationForJ
		})

	case "desc":
		sort.Slice(groupedTechWindows, func(i int, j int) bool {
			durationForI, err := time.ParseDuration(groupedTechWindows[i].AverageDuration)
			if err != nil {
				return false
			}

			durationForJ, err := time.ParseDuration(groupedTechWindows[j].AverageDuration)
			if err != nil {
				return false
			}

			return durationForI > durationForJ
		})

	}

	return groupedTechWindows
}

func AddWindowQuery(objectID int, start string, end string) error {
	err := dbhelper.AddWindowQuery(dbhelper.DB, objectID, start, end)
	if err != nil {
		return fmt.Errorf("add window error: %w", err)
	}

	return nil
}

func AddWindowProc(objectID int, start string, end string) error {
	objects, err := dbhelper.SelectObject(dbhelper.DB)
	if err != nil {
		return fmt.Errorf("select objects error: %w", err)
	}

	if !isContainsObjectWithID(objectID, objects) {
		return errors.New("object with this id is not exists")
	}

	techWindows, err := dbhelper.SelectTechWindows(dbhelper.DB)
	if err != nil {
		return fmt.Errorf("select techt windows error: %w", err)
	}

	filteredTechWindows := filterTechWindowsByObjectID(techWindows, objectID)
	bounds := datastructures.TimeRange{Start: start, End: end}
	if isDurationBoundsOverlapped(bounds, filteredTechWindows) {
		return errors.New("duration overlapped")
	}

	err = dbhelper.AddWindowQuery(dbhelper.DB, objectID, start, end)
	if err != nil {
		return fmt.Errorf("DATABASE add window error: %w", err)
	}

	return nil
}

func isContainsObjectWithID(objectID int, objects []datastructures.Object) bool {
	low := 0
	top := len(objects) - 1

	for low <= top {
		midle := (low + top) / 2
		if objects[midle].ID < objectID {
			low = midle + 1
		} else {
			top = midle - 1
		}
	}

	if low == len(objects) || objects[low].ID != objectID {
		return false
	}

	return true
}

func filterTechWindowsByObjectID(techWindows []datastructures.TechWindow, objectID int) []datastructures.TechWindow {
	result := []datastructures.TechWindow{}
	for _, techWindow := range techWindows {
		if techWindow.IDObject == objectID {
			result = append(result, techWindow)
		}
	}

	return result
}

func isDurationBoundsOverlapped(bounds datastructures.TimeRange, techWindows []datastructures.TechWindow) bool {
	for _, techWindow := range techWindows {
		if bounds.IsOverlapped(techWindow.Duration) {
			return true
		}
	}

	return false
}

func Create(objectsCount int, windowsCount int) error {
	err := dbhelper.CallGenerateFunction(dbhelper.DB, objectsCount, windowsCount)
	if err != nil {
		return fmt.Errorf("call generate function error: %w", err)
	}

	return nil
}

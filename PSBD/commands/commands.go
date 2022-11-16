package commands

import (
	"PSBD/dbhelper"
	"errors"
	"fmt"
)

func GetObjects() []dbhelper.Object {
	return dbhelper.SelectObject(dbhelper.DB)
}

func GetTechWindows() []dbhelper.GroupedTechWindow {
	return dbhelper.SelectGroupedTechWindows(dbhelper.DB)
}

func GetSortedWindowsQuery(sort string, start string, end string) []dbhelper.GroupedTechWindow {
	return dbhelper.SelectSortedTechWindowWithBounds(dbhelper.DB, start, end, sort)
}

func GetSortedWindowsProc(sort string, start string, end string) []dbhelper.GroupedTechWindow {
	techWindows := dbhelper.SelectTechWindows(dbhelper.DB)
	sortedTechWindows := sortTechWindows(techWindows, sort)
	return sortedTechWindows
}

func sortTechWindows(techWindows []dbhelper.TechWindow, sort string) []dbhelper.GroupedTechWindow {
	fmt.Sprint(techWindows)
	result := []dbhelper.GroupedTechWindow{}
	return result
}

func AddWindowQuery(objectID int, start string, end string) error {
	err := dbhelper.AddWindowQuery(dbhelper.DB, objectID, start, end)
	if err != nil {
		return fmt.Errorf("add window error: %w", err)
	}

	return nil
}

func AddWindowProc(objectID int, start string, end string) error {
	return errors.New("not realized")
}

func Create(objectsCount int, windowsCount int) {
	dbhelper.CallGenerateFunction(dbhelper.DB, objectsCount, windowsCount)
}

package dbhelper

import (
	"PSBD/datastructures"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

var DB *sql.DB

const (
	host     = "postgres"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "postgres"
)

type TechWindow struct {
	Id       int       `json:"id"`
	IDObject int       `json:"id_object"`
	Duration TimeRange `json:"duration"`
}

func (tw TechWindow) String() string {
	return fmt.Sprintf("Tech window(id: %d, id_object: %d, duration: %s)", tw.Id, tw.IDObject, tw.Duration)
}

type TimeRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

const timeParseTemplate = "2006-01-02 15:04:05"

func (thisTR TimeRange) IsContains(tr TimeRange) bool {
	thisStart, err := time.Parse(timeParseTemplate, thisTR.Start)
	if err != nil {
		return false
	}

	thisEnd, err := time.Parse(timeParseTemplate, thisTR.End)
	if err != nil {
		return false
	}

	dStart, err := time.Parse(timeParseTemplate, tr.Start)
	if err != nil {
		return false
	}

	dEnd, err := time.Parse(timeParseTemplate, tr.End)
	if err != nil {
		return false
	}

	return (thisStart.Before(dStart) || thisStart.Equal(dStart)) && thisEnd.After(dEnd)
}

const dateAsStringLength = 19

func (thisTR TimeRange) IsOverlapped(tr TimeRange) bool {
	thisStart, err := time.Parse(timeParseTemplate, thisTR.Start[:dateAsStringLength])
	if err != nil {
		return false
	}

	thisEnd, err := time.Parse(timeParseTemplate, thisTR.End[:dateAsStringLength])
	if err != nil {
		return false
	}

	trStart, err := time.Parse(timeParseTemplate, tr.Start[:dateAsStringLength])
	if err != nil {
		return false
	}

	trEnd, err := time.Parse(timeParseTemplate, tr.End[:dateAsStringLength])
	if err != nil {
		return false
	}

	return !(trEnd.Before(thisStart) || trEnd.Equal(thisStart)) &&
		!(trStart.After(thisEnd) || trStart.Equal(thisEnd))
}

func (tr TimeRange) GetDuration() time.Duration {
	start, err := time.Parse(timeParseTemplate, tr.Start)
	if err != nil {
		return time.Duration(0)
	}

	end, err := time.Parse(timeParseTemplate, tr.End)
	if err != nil {
		return time.Duration(0)
	}

	return end.Sub(start)
}

func (tr TimeRange) String() string {
	return fmt.Sprintf("[%s, %s)", tr.Start, tr.End)
}

type GroupedTechWindow struct {
	IDObject        int    `json:"id_object"`
	WindowsCount    int    `json:"windows_count"`
	AverageDuration string `json:"average_duration"`
}

func ConnectToDb() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("Connected!")
	return db
}

func SelectObject(db *sql.DB) []datastructures.Object {
	query := getSelectAllObjectsQuery()
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return convertRowsToObjectSlice(rows)
}

func SelectTechWindows(db *sql.DB) []TechWindow {
	query := getSelectAllTechWindowsQuery()
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return convertRowsToTechWindowSlice(rows)
}

func SelectGroupedTechWindows(db *sql.DB) []GroupedTechWindow {
	query := getSelectAllTechWindowsGroupByObjectQuery()
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return convertRowsToGroupedTechWindowSlice(rows)
}

func SelectSortedTechWindowWithBounds(db *sql.DB, start string, end string, sort string) []GroupedTechWindow {
	query := getSelectSortedTechWindowWithBounds(sort, start, end)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return convertRowsToGroupedTechWindowSlice(rows)
}

func CallGenerateFunction(db *sql.DB, objectsCount int, windowsCount int) {
	query := getGenerateFunctionCallQuery(objectsCount, windowsCount)
	_, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
}

func AddWindowQuery(db *sql.DB, objectID int, start string, end string) error {
	query := getAddWindowQuery(objectID, start, end)
	_, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("query execution error: %w", err)
	}

	return nil
}

func convertRowsToObjectSlice(rows *sql.Rows) []datastructures.Object {
	objects := []datastructures.Object{}

	for rows.Next() {
		o := datastructures.Object{}
		err := rows.Scan(&o.ID, &o.Name, &o.Clock)
		if err != nil {
			fmt.Println(err)
			continue
		}

		objects = append(objects, o)
	}

	return objects
}

func convertRowsToTechWindowSlice(rows *sql.Rows) []TechWindow {
	techWindows := []TechWindow{}

	for rows.Next() {
		tw := TechWindow{}
		durationAsBytes := []uint8{}
		err := rows.Scan(&tw.Id, &tw.IDObject, &durationAsBytes)
		if err != nil {
			fmt.Println(err)
			continue
		}
		durationAsString := string(durationAsBytes)
		tw.Duration = ConvertStringToDuration(durationAsString)
		techWindows = append(techWindows, tw)
	}

	return techWindows
}

func convertRowsToGroupedTechWindowSlice(rows *sql.Rows) []GroupedTechWindow {
	groupedTechWindows := []GroupedTechWindow{}

	for rows.Next() {
		gtw := GroupedTechWindow{}
		durationAsBytes := []uint8{}
		err := rows.Scan(&gtw.IDObject, &gtw.WindowsCount, &durationAsBytes)
		if err != nil {
			fmt.Println(err)
			continue
		}
		gtw.AverageDuration = string(durationAsBytes)
		groupedTechWindows = append(groupedTechWindows, gtw)
	}

	return groupedTechWindows
}

// duration as string pattern: ["YYYY-MM-DD HH:MI:SS","YYYY-MM-DD HH:MI:SS")
func ConvertStringToDuration(durationAsString string) TimeRange {
	trimmedString := strings.Trim(durationAsString, "[\")")
	splittedString := strings.Split(trimmedString, "\",\"")
	return TimeRange{Start: splittedString[0], End: splittedString[1]}
}

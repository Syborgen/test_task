package dbhelper

import (
	"database/sql"
	"errors"
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

type Object struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Clock int    `json:"clock"`
}

func (o Object) String() string {
	return fmt.Sprintf("Object(id: %d, name: %s, clock: %d)", o.ID, o.Name, o.Clock)
}

type TechWindow struct {
	Id       int      `json:"id"`
	IDObject int      `json:"id_object"`
	Duration Duration `json:"duration"`
}

func (tw TechWindow) String() string {
	return fmt.Sprintf("Tech window(id: %d, id_object: %d, duration: %s)", tw.Id, tw.IDObject, tw.Duration)
}

type Duration struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

const timeParseTemplate = "2006-01-02 15:04:05"

func NewDuration(start string, end string) error {
	startTime, err := time.Parse(timeParseTemplate, start)
	if err != nil {
		return fmt.Errorf("Parse error: %w", err)
	}

	endTime, err := time.Parse(timeParseTemplate, end)
	if err != nil {
		return fmt.Errorf("Parse error: %w", err)
	}

	if !startTime.Before(endTime) {
		return errors.New("start must be before end")
	}

	return nil
}

func (this Duration) IsContains(d Duration) bool {
	thisStart, err := time.Parse(timeParseTemplate, this.Start)
	if err != nil {
		return false
	}

	thisEnd, err := time.Parse(timeParseTemplate, this.End)
	if err != nil {
		return false
	}

	dStart, err := time.Parse(timeParseTemplate, d.Start)
	if err != nil {
		return false
	}

	dEnd, err := time.Parse(timeParseTemplate, d.End)
	if err != nil {
		return false
	}

	return (thisStart.Before(dStart) || thisStart.Equal(dStart)) && thisEnd.After(dEnd)
}

const dateAsStringLength = 19

func (this Duration) IsOverlapped(d Duration) bool {
	thisStart, err := time.Parse(timeParseTemplate, this.Start[:dateAsStringLength])
	if err != nil {
		return false
	}

	thisEnd, err := time.Parse(timeParseTemplate, this.End[:dateAsStringLength])
	if err != nil {
		return false
	}

	dStart, err := time.Parse(timeParseTemplate, d.Start[:dateAsStringLength])
	if err != nil {
		return false
	}

	dEnd, err := time.Parse(timeParseTemplate, d.End[:dateAsStringLength])
	if err != nil {
		return false
	}

	return !dEnd.Before(thisStart) && !(dStart.After(thisEnd) || dStart.Equal(thisEnd))
}

func (d Duration) GetDuration() time.Duration {
	start, err := time.Parse(timeParseTemplate, d.Start)
	if err != nil {
		return time.Duration(0)
	}

	end, err := time.Parse(timeParseTemplate, d.End)
	if err != nil {
		return time.Duration(0)
	}

	return end.Sub(start)
}

func (d Duration) String() string {
	return fmt.Sprintf("[%s, %s)", d.Start, d.End)
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

func SelectObject(db *sql.DB) []Object {
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

func convertRowsToObjectSlice(rows *sql.Rows) []Object {
	objects := []Object{}

	for rows.Next() {
		o := Object{}
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
func ConvertStringToDuration(durationAsString string) Duration {
	trimmedString := strings.Trim(durationAsString, "[\")")
	splittedString := strings.Split(trimmedString, "\",\"")
	return Duration{Start: splittedString[0], End: splittedString[1]}
}

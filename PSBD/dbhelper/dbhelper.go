package dbhelper

import (
	"database/sql"
	"fmt"
	"strings"
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
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Clock int    `json:"clock"`
}

func (o Object) String() string {
	return fmt.Sprintf("Object(id: %d, name: %s, clock: %d)", o.Id, o.Name, o.Clock)
}

type TechWindow struct {
	Id       int      `json:"id"`
	IdObject int      `json:"id_object"`
	Duration duration `json:"duration"`
}

func (tw TechWindow) String() string {
	return fmt.Sprintf("Tech window(id: %d, id_object: %d, duration: %s)", tw.Id, tw.IdObject, tw.Duration)
}

type duration struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func (d duration) String() string {
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
		err := rows.Scan(&o.Id, &o.Name, &o.Clock)
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
		err := rows.Scan(&tw.Id, &tw.IdObject, &durationAsBytes)
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
func ConvertStringToDuration(durationAsString string) duration {
	trimmedString := strings.Trim(durationAsString, "[\")")
	splittedString := strings.Split(trimmedString, "\",\"")
	return duration{Start: splittedString[0], End: splittedString[1]}
}

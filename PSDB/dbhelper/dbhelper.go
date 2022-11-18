package dbhelper

import (
	"PSBD/datastructures"
	"database/sql"
	"fmt"
)

var DB *sql.DB

const (
	host     = "postgres"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "postgres"
)

func ConnectToDb() (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("open db error: %w", err)
	}
	fmt.Println("Connected!")
	return db, nil
}

func SelectObject(db *sql.DB) ([]datastructures.Object, error) {
	query := getSelectAllObjectsQuery()
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %w", err)
	}

	return convertRowsToObjectSlice(rows)
}

func SelectTechWindows(db *sql.DB) ([]datastructures.TechWindow, error) {
	query := getSelectAllTechWindowsQuery()
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %w", err)
	}

	return convertRowsToTechWindowSlice(rows)
}

func SelectGroupedTechWindows(db *sql.DB) ([]datastructures.GroupedTechWindow, error) {
	query := getSelectAllTechWindowsGroupByObjectQuery()
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %w", err)
	}

	return convertRowsToGroupedTechWindowSlice(rows)
}

func SelectSortedTechWindowWithBounds(
	db *sql.DB, start string, end string, sort string) ([]datastructures.GroupedTechWindow, error) {
	query := getSelectSortedTechWindowWithBounds(sort, start, end)
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %w", err)
	}

	return convertRowsToGroupedTechWindowSlice(rows)
}

func CallGenerateFunction(db *sql.DB, objectsCount int, windowsCount int) error {
	query := getGenerateFunctionCallQuery(objectsCount, windowsCount)
	_, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("query execution error: %w", err)
	}

	return nil
}

func AddWindowQuery(db *sql.DB, objectID int, start string, end string) error {
	query := getAddWindowQuery(objectID, start, end)
	_, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("query execution error: %w", err)
	}

	return nil
}

func convertRowsToObjectSlice(rows *sql.Rows) ([]datastructures.Object, error) {
	objects := []datastructures.Object{}

	for rows.Next() {
		o := datastructures.Object{}
		err := rows.Scan(&o.ID, &o.Name, &o.Clock)
		if err != nil {
			return nil, fmt.Errorf("scan rows error: %w", err)
		}

		objects = append(objects, o)
	}

	return objects, nil
}

func convertRowsToTechWindowSlice(rows *sql.Rows) ([]datastructures.TechWindow, error) {
	techWindows := []datastructures.TechWindow{}

	for rows.Next() {
		tw := datastructures.TechWindow{}
		durationAsBytes := []uint8{}
		err := rows.Scan(&tw.Id, &tw.IDObject, &durationAsBytes)
		if err != nil {
			return nil, fmt.Errorf("scan rows error: %w", err)
		}
		durationAsString := string(durationAsBytes)
		tw.Duration = datastructures.NewTimeRange(durationAsString)
		techWindows = append(techWindows, tw)
	}

	return techWindows, nil
}

func convertRowsToGroupedTechWindowSlice(rows *sql.Rows) ([]datastructures.GroupedTechWindow, error) {
	groupedTechWindows := []datastructures.GroupedTechWindow{}

	for rows.Next() {
		gtw := datastructures.GroupedTechWindow{}
		durationAsBytes := []uint8{}
		err := rows.Scan(&gtw.IDObject, &gtw.WindowsCount, &durationAsBytes)
		if err != nil {
			return nil, fmt.Errorf("scan rows error: %w", err)
		}
		gtw.AverageDuration = string(durationAsBytes)
		groupedTechWindows = append(groupedTechWindows, gtw)
	}

	return groupedTechWindows, nil
}

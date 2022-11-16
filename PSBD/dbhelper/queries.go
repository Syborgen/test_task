package dbhelper

import "fmt"

func getSelectAllObjectsQuery() string {
	return "select * from objects;"
}

func getSelectAllTechWindowsQuery() string {
	return "select * from tech_windows;"
}

func getSelectAllTechWindowsGroupByObjectQuery() string {
	return `SELECT id_object, count(duration) AS windows_count, avg(get_interval(duration)) AS average_duration 
	FROM tech_windows 
	GROUP BY id_object;`
}

func getSelectSortedTechWindowWithBounds(sort string, start string, end string) string {
	return fmt.Sprintf(
		`SELECT id_object, count(duration) AS windows_count, avg(get_interval(duration)) AS average_duration 
		FROM tech_windows 
		WHERE duration <@ '[%s,%s]'::tsrange 
		GROUP BY id_object 
		ORDER BY average_duration %s;`,
		start, end, sort,
	)
}

func getGenerateFunctionCallQuery(objects int, windows int) string {
	return fmt.Sprintf(
		`SELECT generate(%d, %d);`,
		objects, windows,
	)
}

func getAddWindowQuery(objectID int, start string, end string) string {
	return fmt.Sprintf(
		`INSERT INTO tech_windows(id_object, duration) VALUES(%d, '[%s, %s)');`,
		objectID, start, end,
	)
}

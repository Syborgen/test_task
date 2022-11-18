package datastructures

import "fmt"

type TechWindow struct {
	Id       int       `json:"id"`
	IdObject int       `json:"id_object"`
	Duration TimeRange `json:"duration"`
}

const TechWindowRowFormat = "%5v|%15v|%45v"

func (tw TechWindow) String() string {
	return fmt.Sprintf(TechWindowRowFormat, tw.Id, tw.IdObject, tw.Duration.String())
}

func CreateTechWindowTable(techWindows []TechWindow) string {
	var table string

	table += "Technical windows:\n"
	table += fmt.Sprintf(TechWindowRowFormat+"\n", "id", "object id", "duration")

	for _, techWindow := range techWindows {
		table += techWindow.String() + "\n"
	}

	return table
}

func CreateGroupedTechWindowTable(groupedTechWindows []GroupedTechWindow) string {
	var table string

	table += "Grouped technical windows:\n"
	table += fmt.Sprintf(GroupedTechWindowRowFormat+"\n", "object id", "windows count", "average duration")

	for _, groupedTechWindow := range groupedTechWindows {
		table += groupedTechWindow.String() + "\n"
	}

	return table
}

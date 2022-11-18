package datastructures

import "fmt"

const ObjectRowFormat = "%5v|%40v|%5v"

type Object struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Clock int    `json:"clock"`
}

func (o Object) String() string {
	return fmt.Sprintf(ObjectRowFormat, o.Id, o.Name, o.Clock)
}

func CreateObjectTable(objects []Object) string {
	var table string

	table += "Objects:\n"
	table += fmt.Sprintf(ObjectRowFormat+"\n", "id", "name", "clock")

	for _, object := range objects {
		table += object.String() + "\n"
	}

	return table
}

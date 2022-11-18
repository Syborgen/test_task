package datastructures

import "fmt"

type Object struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Clock int    `json:"clock"`
}

func (o Object) String() string {
	return fmt.Sprintf("Object(id: %d, name: %s, clock: %d)", o.ID, o.Name, o.Clock)
}

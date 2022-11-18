package datastructures

import "fmt"

type TechWindow struct {
	Id       int       `json:"id"`
	IDObject int       `json:"id_object"`
	Duration TimeRange `json:"duration"`
}

func (tw TechWindow) String() string {
	return fmt.Sprintf("Tech window(id: %d, id_object: %d, duration: %s)", tw.Id, tw.IDObject, tw.Duration)
}

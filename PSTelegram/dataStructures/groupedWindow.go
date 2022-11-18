package datastructures

import "fmt"

const GroupedTechWindowRowFormat = "%10v|%15v|%30v"

type GroupedTechWindow struct {
	IDObject        int    `json:"id_object"`
	WindowsCount    int    `json:"windows_count"`
	AverageDuration string `json:"average_duration"`
}

func (gtw GroupedTechWindow) String() string {
	return fmt.Sprintf(GroupedTechWindowRowFormat, gtw.IDObject, gtw.WindowsCount, gtw.AverageDuration)
}

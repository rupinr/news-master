package common

type TimeSlot string

const (
	Morning   TimeSlot = "Morning"
	Afternoon TimeSlot = "Afternoon"
	Evening   TimeSlot = "Evening"
	Night     TimeSlot = "Night"
)

type Test struct {
	Some string
}

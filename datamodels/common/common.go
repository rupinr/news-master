package common

type TimeSlot string

const (
	Morning     TimeSlot = "Morning"     // 6:00 - 12:00
	Afternoon   TimeSlot = "Afternoon"   // 12:00 - 18:00
	Evening     TimeSlot = "Evening"     // 18:00 - 20:00
	Night       TimeSlot = "Night"       // 20:00 - 23:00
	SilentHours TimeSlot = "SilentHours" // 23:00 - 06:00

)

type Test struct {
	Some string
}

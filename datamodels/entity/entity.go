package entity

import (
	"news-master/datamodels/common"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Topic struct {
	gorm.Model
	Name    string
	Visible bool
}
type Site struct {
	gorm.Model
	Url    string
	Active bool
}
type User struct {
	gorm.Model
	Email string
}
type Subscription struct {
	gorm.Model
	UserID                 uint
	User                   User
	Topics                 pq.StringArray `gorm:"type:text[]"`
	Sites                  pq.StringArray `gorm:"type:text[]"`
	SubscriptionScheduleID uint
	SubscriptionSchedule   SubscriptionSchedule
	Confirmed              bool
}
type SubscriptionSchedule struct {
	gorm.Model
	Sunday       bool
	Monday       bool
	Tuesday      bool
	Wednesday    bool
	Thursday     bool
	Friday       bool
	Saturday     bool
	TimeSlotEnum common.TimeSlot
	TimeZone     string
}

type Article struct {
	gorm.Model
	Status       string
	TotalResults int
	Results      []Result
	NextPage     string
}

type Result struct {
	Title       string
	Link        string
	Description string
	Content     string
	ImageURL    string
	Language    string
	Country     []string
	Category    []string
	AIRegion    string
}

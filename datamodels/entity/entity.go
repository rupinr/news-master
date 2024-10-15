package entity

import (
	"news-master/datamodels/common"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Topic struct {
	gorm.Model
	Name    string `gorm:"uniqueIndex:idx_name_visible"`
	Visible bool   `gorm:"index:idx_name_visible"`
}
type Site struct {
	gorm.Model
	Url    string
	Active bool `gorm:"index:"`
}
type User struct {
	gorm.Model
	Email             string `gorm:"uniqueIndex"`
	LoginAttemptCount int
}
type Subscription struct {
	gorm.Model
	UserID                 uint `gorm:"uniqueIndex"`
	User                   User
	Topics                 pq.StringArray `gorm:"type:text[]"`
	Sites                  pq.StringArray `gorm:"type:text[]"`
	SubscriptionScheduleID uint
	SubscriptionSchedule   SubscriptionSchedule
	Confirmed              bool
	LastProcessedAt        time.Time
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
	Title          string
	Link           string
	Description    string
	Content        string
	ImageURL       string
	Language       string
	Country        pq.StringArray `gorm:"type:text[]"`
	Category       pq.StringArray `gorm:"type:text[]"`
	DetectedTopics pq.StringArray `gorm:"type:text[]"`
}

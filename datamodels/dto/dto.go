package dto

import (
	"news-master/datamodels/common"

	"github.com/lib/pq"
)

type Topic struct {
	Name    string `json:"name" binding:"required"`
	Visible bool   `json:"visible"`
}

type TopicUpdate struct {
	Visibility *bool `json:"visibility" binding:"required"`
}

type Site struct {
	Url string `json:"url" binding:"required"`
}

type User struct {
	Email string `json:"email" binding:"required"`
}

type Subscription struct {
	Confirmed                bool                 `json:"confirmed"`
	Email                    string               `json:"email" binding:"required"`
	Topics                   pq.StringArray       `json:"topics" binding:"required"`
	Sites                    pq.StringArray       `json:"sites" binding:"required"`
	SubscriptionScheduleData SubscriptionSchedule `json:"subscriptionSchedule" binding:"required"`
}

type SubscriptionConfirmation struct {
	Confirmed *bool `json:"confirmed" binding:"required"`
}

type DailyFrequency struct {
	Monday    bool `json:"monday" binding:"required"`
	Tuesday   bool `json:"tuesday" binding:"required"`
	Wednesday bool `json:"wednesday" binding:"required"`
	Thursday  bool `json:"thursday" binding:"required"`
	Friday    bool `json:"friday" binding:"required"`
	Saturday  bool `json:"staturday" binding:"required"`
	Sunday    bool `json:"sunday" binding:"required"`
}

type SubscriptionSchedule struct {
	DailyFrequency DailyFrequency  `json:"dailyFrequency" binding:"required"`
	TimeSlot       common.TimeSlot `json:"timeSlot" binding:"required"`
	TimeZone       string          `json:"timezone" binding:"required"`
}

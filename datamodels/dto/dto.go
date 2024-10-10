package dto

import (
	"news-master/datamodels/common"

	"github.com/lib/pq"
)

type Topic struct {
	Name string `form:"name"`
}

type Site struct {
	Url string `form:"url"`
}

type User struct {
	Email string `form:"email"`
}

type Subscription struct {
	Email                    string               `form:"email" json:"email"`
	Topics                   pq.StringArray       `form:"topics" json:"topics"`
	Sites                    pq.StringArray       `form:"sites" json:"sites"`
	SubscriptionScheduleData SubscriptionSchedule `form:"subscriptionSchedule" json:"subscriptionSchedule"`
}

type DailyFrequency struct {
	Monday    bool `form:"monday" json:"monday"`
	Tuesday   bool `form:"tuesday" json:"tuesday"`
	Wednesday bool `form:"wednesday" json:"wednesday"`
	Thursday  bool `form:"thursday" json:"thursday"`
	Friday    bool `form:"friday" json:"friday"`
	Saturday  bool `form:"staturday" json:"staturday"`
	Sunday    bool `form:"sunday" json:"sunday"`
}

type SubscriptionSchedule struct {
	DailyFrequency DailyFrequency  `form:"dailyFrequency" json:"dailyFrequency"`
	TimeSlot       common.TimeSlot `form:"timeSlot" json:"timeSlot"`
	TimeZone       string          `form:"timezone" json:"timezone"`
}

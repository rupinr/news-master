package main

import (
	"news-master/cmd/process"
	"news-master/datamodels/common"
	"news-master/datamodels/entity"
	"time"
)

func main() {
	location, _ := time.LoadLocation("Europe/Berlin")
	process.Notify(time.Date(2024, time.October, 10, 20, 48, 0, 0, location), entity.Subscription{
		User:                 entity.User{Email: "rupinr@gmail.com"},
		SubscriptionSchedule: entity.SubscriptionSchedule{Thursday: true, TimeSlotEnum: common.TimeSlot("Night"), TimeZone: "Europe/Berlin"}})
}

package main

import (
	"news-master/actions"
	"news-master/cmd/process"
	"time"
)

func main() {
	location, _ := time.LoadLocation("Europe/Berlin")
	process.Notify(time.Date(2024, time.October, 10, 17, 0, 0, 0, location), actions.Subscription{
		User:                 actions.User{Email: "rupinr@gmail.com"},
		SubscriptionSchedule: actions.SubscriptionSchedule{Thursday: true, Afternoon: true, TimeZone: "Europe/Berlin"}})
}

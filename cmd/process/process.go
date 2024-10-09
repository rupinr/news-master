package process

import (
	"fmt"
	"news-master/actions"
	"time"
)

func Notify(currentServerTime time.Time, subscription actions.Subscription) {
	location, _ := time.LoadLocation(subscription.SubscriptionSchedule.TimeZone)
	weekdayInLocation := currentServerTime.In(location).Weekday()
	if enabledOnSunday(subscription) && isSundayInLocation(weekdayInLocation) {
		fireNotificationInTimeSlot(currentServerTime, subscription)
	} else if enabledOnMonday(subscription) && isMondayInLocation(weekdayInLocation) {
		fireNotificationInTimeSlot(currentServerTime, subscription)
	} else if enabledOnTuesday(subscription) && isTuesdayInLocation(weekdayInLocation) {
		fireNotificationInTimeSlot(currentServerTime, subscription)
	} else if enabledOnWednesday(subscription) && isWednesdayInLocation(weekdayInLocation) {
		fireNotificationInTimeSlot(currentServerTime, subscription)
	} else if enabledOnThursday(subscription) && isThursdayInLocation(weekdayInLocation) {
		fireNotificationInTimeSlot(currentServerTime, subscription)
	} else if enabledOnFriday(subscription) && isFridayInLocation(weekdayInLocation) {
		fireNotificationInTimeSlot(currentServerTime, subscription)
	} else if enabledOnSaturday(subscription) && isSaturdayInLocation(weekdayInLocation) {
		fireNotificationInTimeSlot(currentServerTime, subscription)
	}
}

func isSundayInLocation(weekday time.Weekday) bool {
	return weekday == time.Sunday
}

func isMondayInLocation(weekday time.Weekday) bool {
	return weekday == time.Monday
}
func isTuesdayInLocation(weekday time.Weekday) bool {
	return weekday == time.Tuesday
}
func isWednesdayInLocation(weekday time.Weekday) bool {
	return weekday == time.Wednesday
}
func isThursdayInLocation(weekday time.Weekday) bool {
	return weekday == time.Thursday
}
func isFridayInLocation(weekday time.Weekday) bool {
	return weekday == time.Friday
}
func isSaturdayInLocation(weekday time.Weekday) bool {
	return weekday == time.Saturday
}

func enabledOnSunday(subscription actions.Subscription) bool {
	return subscription.SubscriptionSchedule.Sunday
}

func enabledOnMonday(subscription actions.Subscription) bool {
	return subscription.SubscriptionSchedule.Monday
}

func enabledOnTuesday(subscription actions.Subscription) bool {
	return subscription.SubscriptionSchedule.Tuesday
}
func enabledOnWednesday(subscription actions.Subscription) bool {
	return subscription.SubscriptionSchedule.Wednesday
}
func enabledOnThursday(subscription actions.Subscription) bool {
	return subscription.SubscriptionSchedule.Thursday
}
func enabledOnFriday(subscription actions.Subscription) bool {
	return subscription.SubscriptionSchedule.Friday
}
func enabledOnSaturday(subscription actions.Subscription) bool {
	return subscription.SubscriptionSchedule.Saturday
}

func fireNotificationInTimeSlot(timeInLocation time.Time, subscription actions.Subscription) {
	if subscription.SubscriptionSchedule.Morning && timeSlotInLocation(timeInLocation, subscription.SubscriptionSchedule) == Morning {
		fmt.Printf("%v can get Morning notifcation\n", subscription.User.Email)
	} else if subscription.SubscriptionSchedule.Afternoon && timeSlotInLocation(timeInLocation, subscription.SubscriptionSchedule) == Afternoon {
		fmt.Printf("%v can get Afternoon notifcation\n", subscription.User.Email)
	} else if subscription.SubscriptionSchedule.Evening && timeSlotInLocation(timeInLocation, subscription.SubscriptionSchedule) == Evening {
		fmt.Printf("%v can get Evening notifcation\n", subscription.User.Email)
	} else if subscription.SubscriptionSchedule.Night && timeSlotInLocation(timeInLocation, subscription.SubscriptionSchedule) == Night {
		fmt.Printf("%v can get Night notifcation\n", subscription.User.Email)
	}
}

func timeSlotInLocation(currentServerTime time.Time, schedule actions.SubscriptionSchedule) TimeSlot {
	location, _ := time.LoadLocation(schedule.TimeZone)
	localTime := currentServerTime.In(location)
	if isMorning(localTime, *location) {
		return Morning
	} else if isAfterNoon(localTime, *location) {
		return Afternoon
	} else if isEvening(localTime, *location) {
		return Evening
	} else {
		return Night
	}
}

func isMorning(localtime time.Time, location time.Location) bool {
	six := timeOf(localtime, 6, &location)
	twelve := timeOf(localtime, 12, &location)
	return localtime.Equal(six) || localtime.After(six) && localtime.Equal(twelve) || localtime.Before(twelve)
}

func isAfterNoon(localtime time.Time, location time.Location) bool {
	twelve := timeOf(localtime, 12, &location)
	eighteen := timeOf(localtime, 18, &location)
	return localtime.After(twelve) && localtime.Equal(eighteen) || localtime.Before(eighteen)
}

func isEvening(localtime time.Time, location time.Location) bool {
	eighteen := timeOf(localtime, 18, &location)
	twenty := timeOf(localtime, 12, &location)
	return localtime.After(eighteen) && localtime.Equal(twenty) || localtime.Before(twenty)
}

func timeOf(localTime time.Time, hour int, location *time.Location) time.Time {
	return time.Date(localTime.Year(), localTime.Month(), localTime.Day(), hour, 0, 0, 0, location)
}

type TimeSlot int

const (
	Morning TimeSlot = iota
	Afternoon
	Evening
	Night
)

package process

import (
	"news-master/datamodels/common"
	"news-master/datamodels/entity"
	"time"
)

func Notify(currentServerTime time.Time, subscription entity.Subscription) bool {
	location, _ := time.LoadLocation(subscription.SubscriptionSchedule.TimeZone)
	weekdayInLocation := currentServerTime.In(location).Weekday()
	if enabledOnSunday(subscription) && isSundayInLocation(weekdayInLocation) {
		return fireNotificationInTimeSlot(currentServerTime, subscription)
	} else if enabledOnMonday(subscription) && isMondayInLocation(weekdayInLocation) {
		return fireNotificationInTimeSlot(currentServerTime, subscription)
	} else if enabledOnTuesday(subscription) && isTuesdayInLocation(weekdayInLocation) {
		return fireNotificationInTimeSlot(currentServerTime, subscription)
	} else if enabledOnWednesday(subscription) && isWednesdayInLocation(weekdayInLocation) {
		return fireNotificationInTimeSlot(currentServerTime, subscription)
	} else if enabledOnThursday(subscription) && isThursdayInLocation(weekdayInLocation) {
		return fireNotificationInTimeSlot(currentServerTime, subscription)
	} else if enabledOnFriday(subscription) && isFridayInLocation(weekdayInLocation) {
		return fireNotificationInTimeSlot(currentServerTime, subscription)
	} else if enabledOnSaturday(subscription) && isSaturdayInLocation(weekdayInLocation) {
		return fireNotificationInTimeSlot(currentServerTime, subscription)
	}
	return false
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

func enabledOnSunday(subscription entity.Subscription) bool {
	return subscription.SubscriptionSchedule.Sunday
}

func enabledOnMonday(subscription entity.Subscription) bool {
	return subscription.SubscriptionSchedule.Monday
}

func enabledOnTuesday(subscription entity.Subscription) bool {
	return subscription.SubscriptionSchedule.Tuesday
}
func enabledOnWednesday(subscription entity.Subscription) bool {
	return subscription.SubscriptionSchedule.Wednesday
}
func enabledOnThursday(subscription entity.Subscription) bool {
	return subscription.SubscriptionSchedule.Thursday
}
func enabledOnFriday(subscription entity.Subscription) bool {
	return subscription.SubscriptionSchedule.Friday
}
func enabledOnSaturday(subscription entity.Subscription) bool {
	return subscription.SubscriptionSchedule.Saturday
}

func fireNotificationInTimeSlot(timeInLocation time.Time, subscription entity.Subscription) bool {
	if subscription.SubscriptionSchedule.TimeSlotEnum == common.Morning && timeSlotInLocation(timeInLocation, subscription.SubscriptionSchedule) == common.Morning {
		return true
	} else if subscription.SubscriptionSchedule.TimeSlotEnum == common.Afternoon && timeSlotInLocation(timeInLocation, subscription.SubscriptionSchedule) == common.Afternoon {
		return true
	} else if subscription.SubscriptionSchedule.TimeSlotEnum == common.Evening && timeSlotInLocation(timeInLocation, subscription.SubscriptionSchedule) == common.Evening {
		return true
	} else if subscription.SubscriptionSchedule.TimeSlotEnum == common.Night && timeSlotInLocation(timeInLocation, subscription.SubscriptionSchedule) == common.Night {
		return true
	}
	return false
}

func timeSlotInLocation(currentServerTime time.Time, schedule entity.SubscriptionSchedule) common.TimeSlot {
	location, _ := time.LoadLocation(schedule.TimeZone)
	localTime := currentServerTime.In(location)
	if isMorning(localTime, *location) {
		return common.Morning
	} else if isAfterNoon(localTime, *location) {
		return common.Afternoon
	} else if isEvening(localTime, *location) {
		return common.Evening
	} else {
		return common.Night
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
	twenty := timeOf(localtime, 20, &location)
	return localtime.After(eighteen) && localtime.Equal(twenty) || localtime.Before(twenty)
}

func timeOf(localTime time.Time, hour int, location *time.Location) time.Time {
	return time.Date(localTime.Year(), localTime.Month(), localTime.Day(), hour, 0, 0, 0, location)
}

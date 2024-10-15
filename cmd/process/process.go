package process

import (
	"fmt"
	"news-master/datamodels/common"
	"news-master/datamodels/entity"
	"time"
)

type CallBack func(uint)

func Notify(currentServerTime *time.Time, subscription *entity.Subscription, updateLastProcessedAt CallBack) bool {
	location, _ := time.LoadLocation(subscription.SubscriptionSchedule.TimeZone)
	weekdayInLocation := currentServerTime.In(location).Weekday()
	if enabledOnSunday(subscription) && isSundayInLocation(weekdayInLocation) ||
		enabledOnMonday(subscription) && isMondayInLocation(weekdayInLocation) ||
		enabledOnTuesday(subscription) && isTuesdayInLocation(weekdayInLocation) ||
		enabledOnWednesday(subscription) && isWednesdayInLocation(weekdayInLocation) ||
		enabledOnThursday(subscription) && isThursdayInLocation(weekdayInLocation) ||
		enabledOnFriday(subscription) && isFridayInLocation(weekdayInLocation) ||
		enabledOnSaturday(subscription) && isSaturdayInLocation(weekdayInLocation) {
		return fireNotificationInTimeSlot(currentServerTime, subscription, updateLastProcessedAt)
	} else {
		return false
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

func enabledOnSunday(subscription *entity.Subscription) bool {
	return subscription.SubscriptionSchedule.Sunday
}

func enabledOnMonday(subscription *entity.Subscription) bool {
	return subscription.SubscriptionSchedule.Monday
}

func enabledOnTuesday(subscription *entity.Subscription) bool {
	return subscription.SubscriptionSchedule.Tuesday
}
func enabledOnWednesday(subscription *entity.Subscription) bool {
	return subscription.SubscriptionSchedule.Wednesday
}
func enabledOnThursday(subscription *entity.Subscription) bool {
	return subscription.SubscriptionSchedule.Thursday
}
func enabledOnFriday(subscription *entity.Subscription) bool {
	return subscription.SubscriptionSchedule.Friday
}
func enabledOnSaturday(subscription *entity.Subscription) bool {
	return subscription.SubscriptionSchedule.Saturday
}

func fireNotificationInTimeSlot(timeInLocation *time.Time, subscription *entity.Subscription, callback CallBack) bool {
	if subscription.SubscriptionSchedule.TimeSlot == common.Morning && timeSlotInLocation(timeInLocation, &subscription.SubscriptionSchedule) == common.Morning ||
		subscription.SubscriptionSchedule.TimeSlot == common.Afternoon && timeSlotInLocation(timeInLocation, &subscription.SubscriptionSchedule) == common.Afternoon ||
		subscription.SubscriptionSchedule.TimeSlot == common.Evening && timeSlotInLocation(timeInLocation, &subscription.SubscriptionSchedule) == common.Evening ||
		subscription.SubscriptionSchedule.TimeSlot == common.Night && timeSlotInLocation(timeInLocation, &subscription.SubscriptionSchedule) == common.Night {
		sendEmail(subscription.User.Email)
		if callback != nil {
			callback(subscription.ID)
		}
		return true
	}
	return false
}

func sendEmail(user string) {

	fmt.Printf("Sending Email to %v\n", user)

}

func timeSlotInLocation(currentServerTime *time.Time, schedule *entity.SubscriptionSchedule) common.TimeSlot {
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

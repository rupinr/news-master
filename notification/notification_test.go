package notification_test

import (
	"news-master/datamodels/common"
	"news-master/datamodels/entity"
	"news-master/notification"
	"testing"
	"time"
)

func TestNotify(t *testing.T) {
	tests := []struct {
		name                 string
		currentServerTime    time.Time
		subscription         entity.Subscription
		expectedNotification bool
	}{
		{
			name:              "Notify on Sunday Morning",
			currentServerTime: time.Date(2023, time.April, 2, 10, 0, 0, 0, time.UTC), // Sunday
			subscription: entity.Subscription{
				SubscriptionSchedule: entity.SubscriptionSchedule{
					TimeZone: "UTC",
					Sunday:   true,
					TimeSlot: common.Morning,
				},
			},
			expectedNotification: true,
		},
		{
			name:              "Notify on Sunday Afternoon",
			currentServerTime: time.Date(2023, time.April, 2, 13, 0, 0, 0, time.UTC), // Sunday
			subscription: entity.Subscription{
				SubscriptionSchedule: entity.SubscriptionSchedule{
					TimeZone: "UTC",
					Sunday:   true,
					TimeSlot: common.Afternoon,
				},
			},
			expectedNotification: true,
		},
		{
			name:              "Notify on Sunday Evening",
			currentServerTime: time.Date(2023, time.April, 2, 19, 0, 0, 0, time.UTC), // Sunday
			subscription: entity.Subscription{
				SubscriptionSchedule: entity.SubscriptionSchedule{
					TimeZone: "UTC",
					Sunday:   true,
					TimeSlot: common.Evening,
				},
			},
			expectedNotification: true,
		},
		{
			name:              "Notify on Sunday Night",
			currentServerTime: time.Date(2023, time.April, 2, 21, 0, 0, 0, time.UTC), // Sunday
			subscription: entity.Subscription{
				SubscriptionSchedule: entity.SubscriptionSchedule{
					TimeZone: "UTC",
					Sunday:   true,
					TimeSlot: common.Night,
				},
			},
			expectedNotification: true,
		},
		{
			name:              "Notify on Sunday Morning in New York",
			currentServerTime: time.Date(2023, time.April, 2, 14, 0, 0, 0, time.UTC), // 10 AM in New York (Eastern Daylight Time)
			subscription: entity.Subscription{
				SubscriptionSchedule: entity.SubscriptionSchedule{
					TimeZone: "America/New_York",
					Sunday:   true,
					TimeSlot: common.Morning,
				},
			},
			expectedNotification: true,
		},
		{
			name:              "Notify on Sunday Afternoon in New York",
			currentServerTime: time.Date(2023, time.April, 2, 17, 0, 0, 0, time.UTC), // 1 PM in New York (Eastern Daylight Time)
			subscription: entity.Subscription{
				SubscriptionSchedule: entity.SubscriptionSchedule{
					TimeZone: "America/New_York",
					Sunday:   true,
					TimeSlot: common.Afternoon,
				},
			},
			expectedNotification: true,
		},
		{
			name:              "Notify on Sunday Evening in New York",
			currentServerTime: time.Date(2023, time.April, 2, 23, 0, 0, 0, time.UTC), // 7 PM in New York (Eastern Daylight Time)
			subscription: entity.Subscription{
				SubscriptionSchedule: entity.SubscriptionSchedule{
					TimeZone: "America/New_York",
					Sunday:   true,
					TimeSlot: common.Evening,
				},
			},
			expectedNotification: true,
		},
		{
			name:              "Notify on Sunday Morning in New York",
			currentServerTime: time.Date(2023, time.April, 2, 14, 0, 0, 0, time.UTC), // 10 AM in New York (Eastern Daylight Time)
			subscription: entity.Subscription{
				SubscriptionSchedule: entity.SubscriptionSchedule{
					TimeZone: "America/New_York",
					Sunday:   true,
					TimeSlot: common.Morning,
				},
			},
			expectedNotification: true,
		},
		{
			name:              "Notify on Sunday Afternoon in New York",
			currentServerTime: time.Date(2023, time.April, 2, 17, 0, 0, 0, time.UTC), // 1 PM in New York (Eastern Daylight Time)
			subscription: entity.Subscription{
				SubscriptionSchedule: entity.SubscriptionSchedule{
					TimeZone: "America/New_York",
					Sunday:   true,
					TimeSlot: common.Afternoon,
				},
			},
			expectedNotification: true,
		},
		{
			name:              "Notify on Sunday Evening in New York",
			currentServerTime: time.Date(2023, time.April, 2, 23, 0, 0, 0, time.UTC), // 7 PM in New York (Eastern Daylight Time)
			subscription: entity.Subscription{
				SubscriptionSchedule: entity.SubscriptionSchedule{
					TimeZone: "America/New_York",
					Sunday:   true,
					TimeSlot: common.Evening,
				},
			},
			expectedNotification: true,
		},
		{
			name:              "No notification on Sunday",
			currentServerTime: time.Date(2023, time.April, 2, 10, 0, 0, 0, time.UTC),
			subscription: entity.Subscription{
				SubscriptionSchedule: entity.SubscriptionSchedule{
					TimeZone: "UTC",
					Sunday:   false,
					TimeSlot: common.Morning,
				},
			},
			expectedNotification: false,
		},
		// No notification outside the specified time slot
		{
			name:              "No notification on Sunday night",
			currentServerTime: time.Date(2023, time.April, 2, 22, 0, 0, 0, time.UTC), // 10 PM
			subscription: entity.Subscription{
				SubscriptionSchedule: entity.SubscriptionSchedule{
					TimeZone: "UTC",
					Sunday:   true,
					TimeSlot: common.Morning,
				},
			},
			expectedNotification: false,
		},

		{
			name:              "No notification on Sunday Morning",
			currentServerTime: time.Date(2023, time.April, 2, 8, 0, 0, 0, time.UTC),
			subscription: entity.Subscription{
				SubscriptionSchedule: entity.SubscriptionSchedule{
					TimeZone: "UTC",
					Sunday:   true,
					TimeSlot: common.Evening,
				},
			},
			expectedNotification: false,
		},
		// Notification not enabled on Monday Evening
		{
			name:              "No notification on Monday Evening",
			currentServerTime: time.Date(2023, time.April, 3, 21, 0, 0, 0, time.UTC),
			subscription: entity.Subscription{
				SubscriptionSchedule: entity.SubscriptionSchedule{
					TimeZone: "UTC",
					Monday:   true,
					TimeSlot: common.Evening,
				},
			},
			expectedNotification: false,
		},
		// Notification not enabled on Tuesday Afternoon
		{
			name:              "No notification on Tuesday Afternoon",
			currentServerTime: time.Date(2023, time.April, 4, 14, 0, 0, 0, time.UTC),
			subscription: entity.Subscription{
				SubscriptionSchedule: entity.SubscriptionSchedule{
					TimeZone: "UTC",
					Tuesday:  true,
					TimeSlot: common.Morning,
				},
			},
			expectedNotification: false,
		},

		// Notification on Morning
		{
			name:              "Notification on Morning",
			currentServerTime: time.Date(2023, time.October, 28, 0, 0, 0, 0, time.UTC),
			subscription: entity.Subscription{
				SubscriptionSchedule: entity.SubscriptionSchedule{
					TimeZone:  "Europe/Berlin",
					Monday:    true,
					Tuesday:   true,
					Wednesday: true,
					Thursday:  true,
					Friday:    true,
					Saturday:  true,
					TimeSlot:  common.Night,
				},
			},
			expectedNotification: false,
		},

		//Silent Hours After 11
		{
			name:              "Silent Hours After 11",
			currentServerTime: time.Date(2024, time.October, 28, 23, 0, 0, 0, time.UTC),
			subscription: entity.Subscription{
				SubscriptionSchedule: entity.SubscriptionSchedule{
					TimeZone: "Europe/Berlin",
					Monday:   true,
					TimeSlot: common.Night,
				},
			},
			expectedNotification: false,
		},

		//No Silent Hours At 22:00
		{
			name:              "No Silent Hours At 22:00 ",
			currentServerTime: time.Date(2024, time.October, 28, 22, 0, 0, 0, time.UTC),
			subscription: entity.Subscription{
				SubscriptionSchedule: entity.SubscriptionSchedule{
					TimeZone: "Europe/Berlin",
					Monday:   true,
					TimeSlot: common.Night,
				},
			},
			expectedNotification: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := notification.IsRightTime(&tt.currentServerTime, &tt.subscription)
			if got != tt.expectedNotification {
				t.Errorf("Failed test %v, Notify() = %v, want %v", tt.name, got, tt.expectedNotification)
			}
		})
	}
}

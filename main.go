package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"news-master/actions"

	"github.com/gin-gonic/gin"
)

type Subscription struct {
	Email  string   `form:"email"`
	Topics []string `form:"topics"`
	Sites  []string `form:"sites"`
}

func main() {
	actions.Migrate()
	r := gin.Default()
	r.POST("/topic", func(c *gin.Context) {
		var topic actions.TopicData
		if c.ShouldBind(&topic) == nil {
			actions.CreateTopic(topic)
		}
		c.String(200, "Success")
	})

	r.POST("/site", func(c *gin.Context) {
		var site actions.SiteData
		if c.ShouldBind(&site) == nil {
			fmt.Println(site)
			actions.CreateSite(site)
		}
		c.String(200, "Success")
	})

	r.POST("/subscribe", func(c *gin.Context) {
		var subscriptionData actions.SubscriptionData
		if c.ShouldBind(&subscriptionData) == nil {
			user := actions.CreateUser(actions.UserData{Email: subscriptionData.Email})

			fmt.Printf("user id is %v\n", user.ID)
			schedule := actions.CreateSubscriptionSchedule(subscriptionData.SubscriptionScheduleData)

			fmt.Printf("schedule id is %v\n", schedule.ID)

			actions.CreateSubscription(subscriptionData, user, schedule)
		}
		c.String(200, "Success")
	})

	r.GET("/subscription", func(c *gin.Context) {
		email := c.Query("email")
		fmt.Println(email)
		sub := actions.GetSubscriptions(email)
		fmt.Printf("sub schedule id is %v\n", sub.SubscriptionSchedule)
		subData := actions.SubscriptionData{
			Email:  email,
			Topics: sub.Topics,
			Sites:  sub.Sites,
			SubscriptionScheduleData: actions.SubscriptionScheduleData{
				DailyFrequency: actions.DailyFrequency{
					Monday:    sub.SubscriptionSchedule.Monday,
					Tuesday:   sub.SubscriptionSchedule.Tuesday,
					Wednesday: sub.SubscriptionSchedule.Wednesday,
					Thursday:  sub.SubscriptionSchedule.Thursday,
					Friday:    sub.SubscriptionSchedule.Friday,
					Saturday:  sub.SubscriptionSchedule.Saturday,
					Sunday:    sub.SubscriptionSchedule.Sunday,
				}, TimeSlot: actions.TimeSlot{
					Morning:   sub.SubscriptionSchedule.Morning,
					Afternoon: sub.SubscriptionSchedule.Afternoon,
					Evening:   sub.SubscriptionSchedule.Evening,
					Night:     sub.SubscriptionSchedule.Night,
				}, TimeZone: sub.SubscriptionSchedule.TimeZone,
			},
		}

		jsonData, _ := json.Marshal(subData)
		c.Data(http.StatusOK, "application/json", jsonData)
	})
	r.Run()
}

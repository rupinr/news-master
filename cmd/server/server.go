package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"news-master/datamodels/dto"
	"news-master/repository"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Subscription struct {
	Email  string   `form:"email"`
	Topics []string `form:"topics"`
	Sites  []string `form:"sites"`
}

func main() {
	repository.Migrate()
	r := gin.Default()
	r.POST("/topic", func(c *gin.Context) {
		var topic dto.Topic
		if c.ShouldBind(&topic) == nil {
			repository.CreateTopic(topic)
		}
		c.String(200, "Success")
	})

	r.POST("/site", func(c *gin.Context) {
		var site dto.Site
		if c.ShouldBind(&site) == nil {
			fmt.Println(site)
			repository.CreateSite(site)
		}
		c.String(200, "Success")
	})

	r.POST("/subscribe", func(c *gin.Context) {
		var subscriptionData dto.Subscription
		if c.ShouldBind(&subscriptionData) == nil {
			user := repository.CreateUser(dto.User{Email: subscriptionData.Email})

			fmt.Printf("user id is %v\n", user.ID)
			schedule := repository.CreateSubscriptionSchedule(subscriptionData.SubscriptionScheduleData)

			fmt.Printf("schedule id is %v\n", schedule.ID)

			sub := repository.CreateSubscription(subscriptionData, user, schedule)
			createdSub := repository.GetSubscriptionByID(int(sub.ID))
			subData := dto.SubscriptionSchedule{DailyFrequency: dto.DailyFrequency{
				Monday:    createdSub.SubscriptionSchedule.Monday,
				Tuesday:   createdSub.SubscriptionSchedule.Tuesday,
				Wednesday: createdSub.SubscriptionSchedule.Wednesday,
				Thursday:  createdSub.SubscriptionSchedule.Thursday,
				Friday:    createdSub.SubscriptionSchedule.Friday,
				Saturday:  createdSub.SubscriptionSchedule.Saturday,
				Sunday:    createdSub.SubscriptionSchedule.Sunday,
			},
				TimeZone: createdSub.SubscriptionSchedule.TimeZone,
				TimeSlot: createdSub.SubscriptionSchedule.TimeSlotEnum}
			s := dto.Subscription{
				Email:                    createdSub.User.Email,
				Topics:                   pq.StringArray(createdSub.Topics),
				Sites:                    pq.StringArray(createdSub.Sites),
				SubscriptionScheduleData: subData,
			}

			jsonData, _ := json.Marshal(s)

			c.Data(http.StatusOK, "application/json", jsonData)
		}
	})

	r.GET("/subscription", func(c *gin.Context) {
		email := c.Query("email")
		fmt.Println(email)
		sub := repository.GetSubscriptionByEmail(email)
		fmt.Printf("sub schedule id is %v\n", sub.SubscriptionSchedule)
		subData := dto.Subscription{
			Email:  sub.User.Email,
			Topics: sub.Topics,
			Sites:  sub.Sites,
			SubscriptionScheduleData: dto.SubscriptionSchedule{
				DailyFrequency: dto.DailyFrequency{
					Monday:    sub.SubscriptionSchedule.Monday,
					Tuesday:   sub.SubscriptionSchedule.Tuesday,
					Wednesday: sub.SubscriptionSchedule.Wednesday,
					Thursday:  sub.SubscriptionSchedule.Thursday,
					Friday:    sub.SubscriptionSchedule.Friday,
					Saturday:  sub.SubscriptionSchedule.Saturday,
					Sunday:    sub.SubscriptionSchedule.Sunday,
				},
				TimeSlot: sub.SubscriptionSchedule.TimeSlotEnum,
				TimeZone: sub.SubscriptionSchedule.TimeZone,
			},
		}

		jsonData, _ := json.Marshal(subData)
		c.Data(http.StatusOK, "application/json", jsonData)
	})
	r.Run()
}

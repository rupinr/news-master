package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"news-master/auth"
	"news-master/datamodels/dto"
	"news-master/env"
	"news-master/repository"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Subscription struct {
	Email  string   `form:"email"`
	Topics []string `form:"topics"`
	Sites  []string `form:"sites"`
}

var loadEnvOnce sync.Once

type Error struct {
	Message string
}

func (err Error) Error() string {
	return err.Message
}

func main() {
	loadEnvOnce.Do(env.LoadEnv)
	repository.Migrate()
	r := gin.Default()

	r.POST("/topic", auth.AuthMiddleware(auth.ValidateAdminToken), func(c *gin.Context) {
		var topic dto.Topic
		if c.ShouldBindJSON(&topic) == nil {
			repository.CreateTopic(topic)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		}
	})

	r.PUT("/topic/:topic", auth.AuthMiddleware(auth.ValidateAdminToken), func(c *gin.Context) {
		topicName := c.Param("topic")
		var update dto.TopicUpdate
		if c.ShouldBindJSON(&update) == nil {
			err := repository.UpdateTopic(topicName, *update.Visibility)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		}
	})

	r.POST("/site", auth.AuthMiddleware(auth.ValidateAdminToken), func(c *gin.Context) {
		var site dto.Site
		if c.ShouldBindJSON(&site) == nil {
			repository.CreateSite(site)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		}
	})

	//TODO add confirmation for subscription.. first save with some unverfied status.
	//Update to verified only when clicking on email..
	r.POST("/subscribe", func(c *gin.Context) {
		var subscriptionData dto.Subscription
		if c.ShouldBindJSON(&subscriptionData) == nil {
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
		} else {
			c.JSON(400, gin.H{"error": "Invalid request"})
		}
	})

	//implement confirm endpoint which validates token and changes verified status to true for subscription
	r.POST("/confirm", auth.AuthMiddleware(auth.ValidateSubscriberToken), func(c *gin.Context) {
		var confirmation dto.SubscriptionConfirmation
		if c.ShouldBindJSON(&confirmation) == nil {
			repository.UpdateSubscriptionConfirmation(confirmation.Email, *confirmation.Confirmed)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		}

	})

	//TODO check token from email param... token is only valid lifelong
	r.GET("/subscription", func(c *gin.Context) {
		email := c.Query("email")
		if email == "" {
			c.JSON(404, gin.H{"error": "Invalid request"})
		} else {
			sub, err := repository.GetSubscriptionByEmail(email)

			if err == nil {
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
			} else {
				c.JSON(404, gin.H{"error": "Invalid request"})
			}
		}

	})
	r.Run()
}

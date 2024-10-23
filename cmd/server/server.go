package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"news-master/app"
	"news-master/auth"
	custom "news-master/custom_validators"
	"news-master/datamodels/dto"
	"news-master/repository"
	"news-master/service"
	"news-master/startup"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Subscription struct {
	Email  string   `form:"email"`
	Topics []string `form:"topics"`
	Sites  []string `form:"sites"`
}

type Error struct {
	Message string
}

func (err Error) Error() string {
	return err.Message
}

func main() {
	startup.Init()
	r := gin.Default()
	validate, e := custom.InitValidator()
	if e != nil {
		fmt.Printf("Error initializing validator: %v\n", e)
		return
	}

	config := cors.New(cors.Config{
		AllowOrigins: []string{app.Config.AllowOrigin},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	})

	r.Use(config)
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

	r.GET("/sites", func(c *gin.Context) {
		sites := repository.GetActiveSites()
		var siteData []dto.Site
		for _, site := range sites {
			siteData = append(siteData, dto.Site{Url: site.Url, Name: site.Name})
		}
		jsonData, _ := json.Marshal(siteData)
		c.Data(http.StatusOK, "application/json", jsonData)
	})

	r.POST("/user", func(c *gin.Context) {
		var user dto.User
		if c.ShouldBindJSON(&user) == nil {
			if err := validate.Struct(user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Email"})
				return
			}
			_, err := service.CreateUserAndTriggerLoginEmail(user)
			if err != nil {
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "Max attempt reached."})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		}
	})

	r.POST("/subscribe", auth.AuthMiddleware(auth.ValidateSubscriberToken), func(c *gin.Context) {
		var subscriptionData dto.Subscription
		err := c.ShouldBindJSON(&subscriptionData)

		fmt.Println(err)
		if err == nil {
			cUser := auth.User(c)

			user := repository.GetUser(dto.User{Email: cUser.Email})

			fmt.Printf("Sitesssss %v\n", subscriptionData.Sites)

			schedule := repository.CreateSubscriptionSchedule(subscriptionData.SubscriptionScheduleData)

			fmt.Printf("schedule from subcribe %v\n", schedule)

			sub := repository.CreateSubscription(user, subscriptionData.Sites, schedule.ID)

			createdSub := repository.GetSubscriptionByID(int(sub.ID))

			subData := dto.SubscriptionSchedule{DailyFrequency: dto.DailyFrequency{
				Monday:    &createdSub.SubscriptionSchedule.Monday,
				Tuesday:   &createdSub.SubscriptionSchedule.Tuesday,
				Wednesday: &createdSub.SubscriptionSchedule.Wednesday,
				Thursday:  &createdSub.SubscriptionSchedule.Thursday,
				Friday:    &createdSub.SubscriptionSchedule.Friday,
				Saturday:  &createdSub.SubscriptionSchedule.Saturday,
				Sunday:    &createdSub.SubscriptionSchedule.Sunday,
			},
				TimeZone: createdSub.SubscriptionSchedule.TimeZone,
				TimeSlot: createdSub.SubscriptionSchedule.TimeSlot,
			}

			s := dto.Subscription{
				Sites:                    sub.Sites,
				SubscriptionScheduleData: subData,
				Confirmed:                subscriptionData.Confirmed,
			}

			jsonData, _ := json.Marshal(s)

			c.Data(http.StatusOK, "application/json", jsonData)
		} else {
			c.JSON(400, gin.H{"error": "Invalid request"})
		}
	})

	r.GET("/subscription", auth.AuthMiddleware(auth.ValidateSubscriberToken), func(c *gin.Context) {
		email := auth.User(c).Email
		if email == "" {
			c.JSON(404, gin.H{"error": "Invalid request"})
		} else {
			sub, err := repository.GetSubscriptionByEmail(email)

			fmt.Printf("sub %v\n", sub)

			if err == nil {
				subData := dto.Subscription{
					Sites:     sub.Sites,
					Confirmed: sub.Confirmed,
					SubscriptionScheduleData: dto.SubscriptionSchedule{
						DailyFrequency: dto.DailyFrequency{
							Monday:    &sub.SubscriptionSchedule.Monday,
							Tuesday:   &sub.SubscriptionSchedule.Tuesday,
							Wednesday: &sub.SubscriptionSchedule.Wednesday,
							Thursday:  &sub.SubscriptionSchedule.Thursday,
							Friday:    &sub.SubscriptionSchedule.Friday,
							Saturday:  &sub.SubscriptionSchedule.Saturday,
							Sunday:    &sub.SubscriptionSchedule.Sunday,
						},
						TimeSlot: sub.SubscriptionSchedule.TimeSlot,
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

	r.POST("/feedback", func(c *gin.Context) {
		var feedback dto.Feedback
		if c.ShouldBindJSON(&feedback) == nil {
			service.CreateFeedBackAndTriggerAdminEmail(feedback)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		}
	})
	r.Run()
}

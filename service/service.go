package service

import (
	"errors"
	"fmt"
	"news-master/app"
	"news-master/auth"
	"news-master/datamodels/common"
	"news-master/datamodels/dto"
	"news-master/datamodels/entity"
	"news-master/email"
	"news-master/helper"
	"news-master/repository"
	"strconv"
)

func CreateUserAndTriggerLoginEmail(user dto.User) (entity.User, error) {
	createdUser := repository.CreateUser(user)
	maxLoginAttempt, _ := strconv.Atoi(app.Config.MaxLoginAttempt)
	if createdUser.LoginAttemptCount < maxLoginAttempt {
		repository.IncrementAndGetLoginAttempt(user)
		token, _ := auth.SubsriberToken(createdUser.ID, user.Email, 24)
		defaultValue := true
		subscriptionSchedule := repository.CreateSubscriptionSchedule(
			dto.SubscriptionSchedule{
				DailyFrequency: dto.DailyFrequency{
					Monday:    &defaultValue,
					Tuesday:   &defaultValue,
					Wednesday: &defaultValue,
					Thursday:  &defaultValue,
					Friday:    &defaultValue,
					Saturday:  &defaultValue,
					Sunday:    &defaultValue,
				},
				TimeSlot: common.Morning,
			},
		)
		repository.CreateSubscription(createdUser, []string{}, subscriptionSchedule.ID)
		go email.SendSesEmail(
			createdUser.Email,
			"Activate Your QuickBrew Subscription Now!",
			"",
			fmt.Sprintf("If the above link doesn't work, please copy and paste the following URL into your browser: %s", helper.PreAuthLink(token)),
		)
		return createdUser, nil
	} else {
		return createdUser, errors.New("max Login attempt reached")
	}

}

package service

import (
	"errors"
	"news-master/auth"
	"news-master/datamodels/common"
	"news-master/datamodels/dto"
	"news-master/datamodels/entity"
	"news-master/email"
	"news-master/repository"
	"os"
	"strconv"
)

func CreateUserAndTriggerLoginEmail(user dto.User) (entity.User, error) {
	createdUser := repository.CreateUser(user)
	maxLoginAttempt, _ := strconv.Atoi(os.Getenv("MAX_LOGIN_ATTEMPT"))
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

		go email.SendEmail(createdUser.Email, token, "activate your email")
		return createdUser, nil
	} else {
		return createdUser, errors.New("max Login attempt reached")
	}

}

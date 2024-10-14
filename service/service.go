package service

import (
	"errors"
	"news-master/auth"
	"news-master/datamodels/dto"
	"news-master/datamodels/entity"
	"news-master/email"
	"news-master/repository"
	"os"
	"strconv"
)

func FirstOrCreateSubscription(
	subscriptionData dto.Subscription,
	user entity.User,
	subscriptionSchedule entity.SubscriptionSchedule,
) entity.Subscription {

	sub := repository.CreateSubscription(subscriptionData, user, subscriptionSchedule)
	return sub

}

func CreateUserAndTriggerLoginEmail(user dto.User) (entity.User, error) {
	createdUser := repository.CreateUser(user)
	maxLoginAttempt, _ := strconv.Atoi(os.Getenv("MAX_LOGIN_ATTEMPT"))
	if createdUser.LoginAttemptCount < maxLoginAttempt {
		repository.IncrementAndGetLoginAttempt(user)
		token, _ := auth.SubsriberToken(createdUser.ID, user.Email, 24)
		go email.SendEmail(createdUser.Email, token, "activate your email")
		return createdUser, nil
	} else {
		return createdUser, errors.New("max Login attempt reached")
	}

}

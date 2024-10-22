package service

import (
	"errors"
	"log/slog"
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
		emailData := email.EmailData{ActivationLink: helper.PreAuthLink(token)}
		htmlEmail, htmlErr := email.GenerateHTML(emailData)
		textEmail, txtErr := email.GenerateText(emailData)
		if htmlErr == nil && txtErr == nil {
			go email.SendEmail(
				createdUser.Email,
				"Activate Your QuickBrew Subscription Now!",
				htmlEmail,
				textEmail,
			)
		} else {
			slog.Error("Error in email template", htmlErr, txtErr)
		}
		return createdUser, nil
	} else {
		return createdUser, errors.New("max Login attempt reached")
	}

}

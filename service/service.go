package service

import (
	"errors"
	"news-master/app"
	"news-master/auth"
	"news-master/datamodels/common"
	"news-master/datamodels/dto"
	"news-master/datamodels/entity"
	"news-master/email"
	"news-master/helper"
	"news-master/logger"
	"news-master/repository"
	"strconv"
)

func CreateUserAndTriggerLoginEmail(user dto.User) (entity.User, error) {
	createdUser, isNewUser, createErr := repository.CreateUser(user)
	maxLoginAttempt, _ := strconv.Atoi(app.Config.MaxLoginAttempt)
	if createdUser.LoginAttemptCount < maxLoginAttempt {
		repository.IncrementAndGetLoginAttempt(user)
		var subject string
		var html string
		var htmlErr error
		token, _ := auth.SubscriberToken(createdUser.ID, user.Email, 24)
		emailData := email.EmailData{ActivationLink: helper.PreAuthLink(token)}
		textEmail, txtErr := email.GenerateText(emailData)
		if isNewUser {
			createSubscriptionSchedule(createdUser)
			subject = "Activate Your QuickBrew Subscription Now!"
			html, htmlErr = email.GenerateRegistrationHTML(emailData)
		} else {
			subject = "Manage Your QuickBrew Subscription"
			html, htmlErr = email.GenerateUpdateHTML(emailData)
		}
		if htmlErr == nil && txtErr == nil && createErr == nil {
			go email.SendEmail(
				createdUser.Email,
				subject,
				html,
				textEmail,
			)
		} else {
			logger.Log.Error("Error in email template", "htmlError", htmlErr.Error(), "txtErr", txtErr.Error())
		}
		logger.Log.Debug("User got registered")
		return createdUser, nil
	} else {
		return createdUser, errors.New("max Login attempt reached")
	}

}

func createSubscriptionSchedule(user entity.User) {
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
			TimeZone: "Europe/Berlin",
		},
	)
	repository.CreateSubscription(user, []string{}, subscriptionSchedule.ID, false)
}

func CreateFeedBackAndTriggerAdminEmail(feedback dto.Feedback) {
	createdFeedback, err := repository.CreateFeedBack(feedback)
	if err == nil {
		go email.SendEmail(app.Config.AdminEmail, "You've got feedback", createdFeedback.Content, createdFeedback.Content)
	}
}

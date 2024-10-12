package service

import (
	"fmt"
	"news-master/auth"
	"news-master/datamodels/dto"
	"news-master/datamodels/entity"
	"news-master/repository"
)

func FirstOrCreateSubscription(
	subscriptionData dto.Subscription,
	user entity.User,
	subscriptionSchedule entity.SubscriptionSchedule,
) entity.Subscription {

	sub := repository.CreateSubscription(subscriptionData, user, subscriptionSchedule)

	fmt.Printf("C Status %v\n", sub.Confirmed)
	if !sub.Confirmed {
		token, err := auth.SubsriberToken(int(user.ID), user.Email, 24)
		fmt.Printf("token error %v\n", err)
		if err == nil {
			go fmt.Printf("Send Email.. inbackground. JWT token is %s", token)
		}
	} else {
		fmt.Println("Already confirmed sub, not need sent email")
	}
	return sub
}

package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"news-master/app"
	"news-master/auth"
	"news-master/datamodels/dto"
	"news-master/email"
	"news-master/helper"
	"news-master/logger"
	"news-master/notification"
	"news-master/repository"
	"time"
)

func FetchNewsTask() {
	sites := repository.GetActiveSites()
	if len(sites) == 0 {
		logger.Log.Warn("No active sites")
	}

	apiRequest, _ := url.Parse(app.Config.NewsDataApiUrl)
	apiRequest.Path = "/api/1/latest"
	params := url.Values{}
	params.Add("apikey", app.Config.NewsDataApiKey)

	for _, site := range sites {
		params.Set("domainurl", site.Url)
		params.Set("language", site.Language)
		apiRequest.RawQuery = params.Encode()

		logger.Log.Debug("Fetching", "url", apiRequest.String())

		resp, apiErr := http.Get(apiRequest.String())

		if apiErr != nil {
			logger.Log.Error("Error fetching News API", "error", apiErr.Error())
			continue
		} else {
			defer resp.Body.Close()
			var response dto.NewsdataApiResponse
			body, readErr := io.ReadAll(resp.Body)

			if readErr != nil {
				logger.Log.Error("Error Reading response from News API", "error", readErr.Error())
				continue
			}

			unmarshalErr := json.Unmarshal(body, &response)

			if unmarshalErr != nil {
				logger.Log.Error("Error Processing response from News API", "error", unmarshalErr.Error())
				continue
			}

			logger.Log.Debug("Found ariticles", "count", len(response.Results))

			for _, result := range response.Results {
				repository.CreateResult(result)
			}
		}
	}
}

func SendNewsletter() {
	subscriptions := repository.GetSubscriptionsToProcess()

	logger.Log.Debug(fmt.Sprintf("Processing %v number of subscriptions", len(subscriptions)))

	for _, subscription := range subscriptions {

		logger.Log.Debug("Number of sites in subscription", "number", len(subscription.Sites), "sub_id", subscription.ID)
		time := time.Now()
		canSendEmail := notification.IsRightTime(&time, &subscription)
		if canSendEmail {
			articles := repository.GetArticlesAfterLastProcessedTime(subscription.LastProcessedAt, subscription.Sites)

			if len(articles) == 0 {
				logger.Log.Debug("No articles found for the subscription, not sending email for", "sub_id", subscription.ID)
				continue
			}

			token, tokenErr := auth.SubscriberToken(subscription.UserID, subscription.User.Email, 24*30) //30 days

			if tokenErr != nil {
				logger.Log.Error("Error generating token for email", "error", tokenErr.Error())
				continue
			}

			html, htmlErr := email.NewsLetterHTML(
				dto.NewsletterData{
					Articles:               articles,
					ManageSubscriptionLink: helper.PreAuthLink(token),
					AboutLink:              helper.AboutLinkLink(),
					PrivacyLink:            helper.PrivacyLink(),
				})

			if htmlErr != nil {
				logger.Log.Error("Error generating HTML", "error", htmlErr.Error())
				continue
			}
			emailError := email.SendEmail(
				subscription.User.Email,
				"Your daily newsletter",
				html,
				"")

			if emailError != nil {
				logger.Log.Error(fmt.Sprintf("Error sending email %v", emailError.Error()))
			} else {
				repository.SetLastProcessedAt(subscription.ID)
				logger.Log.Debug("Sent email for subscription", "sub_id", subscription.ID)
			}
		} else {
			logger.Log.Debug("The subscription is not eligible to recieve email at this slot", "sub_id", subscription.ID)
			continue
		}
	}
}

func CleanUp() {
	deletedRowCount, err := repository.DeleteOldArticlesFrom(time.Now().AddDate(0, 0, -8))
	if err == nil {
		logger.Log.Info("Cleanup Completed successfully", "count", deletedRowCount)
	} else {
		logger.Log.Error("Cleanup Has errors", "error", err.Error())
	}
}

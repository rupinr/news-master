package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"news-master/app"
	"news-master/auth"
	notification "news-master/cmd/process"
	"news-master/datamodels/dto"
	"news-master/email"
	"news-master/helper"
	"news-master/logger"
	"news-master/repository"
	"time"
)

func FetchNewsTask() {
	sites := repository.GetActiveSites()
	if len(sites) == 0 {
		logger.Log.Warn("No active sites")
	}

	for _, site := range sites {

		resp, apiErr := http.Get(fmt.Sprintf("%s/api/1/latest?apikey=%s&domainurl=%s", app.Config.NewsDataApiUrl, app.Config.NewsDataApiKey, site.Url))
		logger.Log.Debug(fmt.Sprintf("Site: %v", site))

		if apiErr != nil {
			logger.Log.Error(fmt.Sprintf("Error fetching News API: %v", apiErr.Error()))
			return
		}
		defer resp.Body.Close()
		var response dto.NewsdataApiResponse
		body, readErr := io.ReadAll(resp.Body)

		if readErr != nil {
			logger.Log.Error(fmt.Sprintf("Error Reading response from News API: %v", readErr.Error()))
			return
		}

		for i := range response.Results {
			response.Results[i].SourceUrl = getDomain(response.Results[i].SourceUrl)
		}

		unmarshalErr := json.Unmarshal(body, &response)

		if unmarshalErr != nil {
			logger.Log.Error(fmt.Sprintf("Error Processing response from News API: %v", unmarshalErr.Error()))
			return
		}

		logger.Log.Debug(fmt.Sprintf("Found %v ariticles", len(response.Results)))

		for _, result := range response.Results {
			repository.CreateResult(result)
		}

	}
}

func getDomain(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Sprintf("INVALID_RAW_URL_%v", rawURL)
	}
	return parsedURL.Host
}

func SendNewsletter() {
	subscriptions := repository.GetSubscriptionsToProcess()

	logger.Log.Debug(fmt.Sprintf("Processing %v number of subscriptions", len(subscriptions)))

	for _, subscription := range subscriptions {

		time := time.Now()
		canSendEmail := notification.IsRightTime(&time, &subscription)
		if canSendEmail {
			articles := repository.GetArticlesAfterLastProcessedTime(subscription.LastProcessedAt, subscription.Sites)

			if len(articles) == 0 {
				logger.Log.Debug("No articles found for the subscription, not sending email")
				continue
			}

			token, tokenErr := auth.SubsriberToken(subscription.UserID, subscription.User.Email, 24)

			if tokenErr != nil {
				logger.Log.Error(fmt.Sprintf("Error generating token for email %v", tokenErr.Error()))
				continue
			}

			html, htmlErr := email.GenerateNewsLetterHTML(dto.NewsletterData{
				Articles:               articles,
				ManageSubscriptionLink: helper.PreAuthLink(token),
				AboutLink:              helper.AboutLinkLink(),
				PrivacyLink:            helper.PrivacyLink(),
			})

			if htmlErr != nil {
				logger.Log.Error(fmt.Sprintf("Error generating HTML for email %v", htmlErr.Error()))
				continue
			}
			emailError := email.SendEmail(
				subscription.User.Email,
				"Your daily newsletter",
				html,
				"")

			if emailError != nil {
				logger.Log.Error(fmt.Sprintf("Unable to send email to %v", emailError.Error()))
			} else {
				logger.Log.Debug(fmt.Sprintf("Setting last processed time stamp for subscription with ID %v", subscription.ID))
				repository.SetLastProcessedAt(subscription.ID)
			}
		} else {
			continue
		}
	}
}

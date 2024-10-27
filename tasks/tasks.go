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

	apiRequest, _ := url.Parse(app.Config.NewsDataApiUrl)
	apiRequest.Path = "/api/1/latest"
	params := url.Values{}
	params.Add("apikey", app.Config.NewsDataApiKey)

	for _, site := range sites {
		params.Set("domainurl", site.Url)
		apiRequest.RawQuery = params.Encode()

		logger.Log.Debug("Fetching", "url", apiRequest.String())

		resp, apiErr := http.Get(apiRequest.String())

		if apiErr != nil {
			logger.Log.Error("Error fetching News API", "error", apiErr.Error())
			return
		}
		defer resp.Body.Close()
		var response dto.NewsdataApiResponse
		body, readErr := io.ReadAll(resp.Body)

		if readErr != nil {
			logger.Log.Error("Error Reading response from News API", "error", readErr.Error())
			return
		}

		for i := range response.Results {
			response.Results[i].SourceUrl = getDomain(response.Results[i].SourceUrl)
		}

		unmarshalErr := json.Unmarshal(body, &response)

		if unmarshalErr != nil {
			logger.Log.Error("Error Processing response from News API", "error", unmarshalErr.Error())
			return
		}

		logger.Log.Debug("Found ariticles", "count", len(response.Results))

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

	for idx, subscription := range subscriptions {

		logger.Log.Debug("Processing Item", "index", idx)

		time := time.Now()
		canSendEmail := notification.IsRightTime(&time, &subscription)
		if canSendEmail {
			articles := repository.GetArticlesAfterLastProcessedTime(subscription.LastProcessedAt, subscription.Sites)

			if len(articles) == 0 {
				logger.Log.Debug("No articles found for the subscription, not sending email")
				continue
			}

			token, tokenErr := auth.SubscriberToken(subscription.UserID, subscription.User.Email, 24*7)

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

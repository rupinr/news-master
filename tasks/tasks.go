package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"news-master/app"
	"news-master/auth"
	"news-master/cmd/process"
	"news-master/datamodels/dto"
	"news-master/email"
	"news-master/helper"
	"news-master/repository"
	"time"
)

func FetchNewsTask() {
	sites := repository.GetActiveSites()
	if len(sites) == 0 {
		slog.Warn("No active sites")
	}

	for _, site := range sites {

		resp, err := http.Get(fmt.Sprintf("%s/api/1/latest?apikey=%s&domainurl=%s", app.Config.NewsDataApiUrl, app.Config.NewsDataApiKey, site.Url))
		fmt.Printf("Site: %v", site)

		if err != nil {
			// handle error
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		var response dto.NewsdataApiResponse
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			// handle error
			fmt.Println("Error:", err)
			return
		}

		errUnmarshal := json.Unmarshal(body, &response)

		if errUnmarshal != nil {
			// handle error
			fmt.Println("Error:", err)
			return
		}

		for _, result := range response.Results {
			repository.CreateResult(result)
		}

	}
}

func SendNewsletter() {
	subscriptions := repository.GetSubscriptionsToProcess()

	if len(subscriptions) == 0 {
		slog.Warn("No confirmed subscrriptions exist")
	}
	for _, subscription := range subscriptions {

		//TODO add logs, when there is no data..
		if subscription.Confirmed {
			time := time.Now()
			fmt.Printf("Last processed at %v\n", subscription.LastProcessedAt)
			articles := repository.GetArticlesFrom(subscription.LastProcessedAt, subscription.Sites)

			//Fix, do not send, it procedded today....
			process.Notify(&time, &subscription, repository.SetLastProcessedAt)

			token, tokenErr := auth.SubsriberToken(subscription.UserID, subscription.User.Email, 24)

			html, err := email.GenerateNewsLetterHTML(dto.NewsletterData{
				Articles:               articles,
				ManageSubscriptionLink: helper.PreAuthLink(token),
				AboutLink:              helper.AboutLinkLink(),
				PrivacyLink:            helper.PrivacyLink(),
			})

			if err == nil && tokenErr == nil {
				email.SendEmail(
					subscription.User.Email,
					"Your daily newsletter",
					html,
					"")
			} else {
				slog.Error("Error sending email", err.Error(), tokenErr.Error())
			}
		}
	}
}

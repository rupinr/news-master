package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"news-master/app"
	"news-master/cmd/process"
	"news-master/datamodels/dto"
	"news-master/repository"
	"time"

	"golang.org/x/exp/rand"
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
	subscriptions := repository.GetAllSubscriptions()

	if len(subscriptions) == 0 {
		slog.Warn("No confirmed subscrriptions exist")
	}
	for _, subscription := range subscriptions {

		//TODO add logs, when there is no data..
		if subscription.Confirmed {
			time := time.Now()
			fmt.Printf("Last processed at %v\n", subscription.LastProcessedAt)
			articles := repository.GetArticlesFrom(subscription.LastProcessedAt, subscription.Sites)
			process.Notify(&time, &subscription, repository.SetLastProcessedAt)
			fmt.Printf("Subscription for %v\n", subscription.User.Email)
			fmt.Printf("SubscriptionTopic for %v\n", subscription.Topics)
			randomIndex := rand.Intn(len(articles))

			fmt.Printf("Articles %v\n", articles[randomIndex].Title)
		}
	}
}

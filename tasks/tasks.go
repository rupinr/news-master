package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"news-master/cmd/process"
	"news-master/datamodels/dto"
	"news-master/repository"
	"os"
	"time"
)

func FetchNewsTask() {
	sites := repository.GetActiveSites()
	for _, site := range sites {
		resp, err := http.Get(fmt.Sprintf("%s/api/1/latest?apikey=%s&domainurl=%s", os.Getenv("NEWS_DATA_API_URL"), os.Getenv("NEWS_DATA_API_KEY"), site.Url))
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
	for _, subscription := range repository.GetAllSubscriptions() {
		if subscription.Confirmed {
			process.Notify(time.Now(), subscription)
			fmt.Printf("Subscription for %v\n", subscription.User.Email)
			fmt.Printf("SubscriptionTopic for %v\n", subscription.Topics)
		}
	}
}

package main

import (
	"fmt"
	"io"
	"net/http"
	"news-master/cmd/process"
	"news-master/repository"
	"news-master/startup"
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func main() {
	startup.Init()
	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
	}

	// add a job to the scheduler

	subscriptionJob, subscriptionJoberr := s.NewJob(
		gocron.CronJob(
			"0 * * * *", false,
		),
		gocron.NewTask(
			func() {
				for _, subscription := range repository.GetAllSubscriptions() {
					if subscription.Confirmed {
						process.Notify(time.Now(), subscription)
						fmt.Printf("Subscription for %v\n", subscription.User.Email)
						fmt.Printf("SubscriptionTopic for %v\n", subscription.Topics)
					}
				}
			},
		),
	)

	newsFetchJob, newsFetchJobErr := s.NewJob(
		gocron.CronJob(
			"* * * * * *", true,
		),
		gocron.NewTask(
			func() {

				sites := repository.GetActiveSites()
				for _, site := range sites {
					resp, err := http.Get(fmt.Sprintf("%s/api/1/latest?apikey=%s&domainurl=%s", os.Getenv("NEWS_DATA_API_URL"), os.Getenv("NEWS_DATA_API_KEY"), site.Url))
					if err != nil {
						// handle error
						fmt.Println("Error:", err)
						return
					}
					defer resp.Body.Close()

					body, err := io.ReadAll(resp.Body)
					if err != nil {
						// handle error
						fmt.Println("Error:", err)
						return
					}
					fmt.Println(string(body))

				}

			},
		),
	)
	if subscriptionJoberr != nil || newsFetchJobErr != nil {
		fmt.Println(subscriptionJoberr)
		fmt.Println(newsFetchJobErr)

		panic("Error in scheduler")
	}
	// each job has a unique id
	fmt.Println(subscriptionJob.ID())
	fmt.Println(newsFetchJob.ID())

	// start the scheduler
	s.Start()

	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute):
	}

	// when you're done, shut it down
	err = s.Shutdown()
	if err != nil {
		fmt.Println("errorrrss")
	}
}

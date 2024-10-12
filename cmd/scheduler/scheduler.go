package main

import (
	"fmt"
	"news-master/cmd/process"
	"news-master/repository"
	"news-master/startup"
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

	j, err := s.NewJob(
		gocron.CronJob(
			"* * * * * *", true,
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

	if err != nil {
		panic("Error in scheduler")
	}
	// each job has a unique id
	fmt.Println(j.ID())

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

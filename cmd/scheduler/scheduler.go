package main

import (
	"fmt"
	"news-master/app"
	"news-master/startup"
	"news-master/tasks"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-co-op/gocron/v2"
)

func main() {
	startup.Init()
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		panic("Unable to start scheduler....")
	}
	_, subscriptionJoberr := scheduler.NewJob(
		gocron.CronJob(app.Config.SubscriptionMailCron, true),
		gocron.NewTask(
			tasks.SendNewsletter,
		),
	)

	_, newsFetchJobErr := scheduler.NewJob(
		gocron.CronJob(app.Config.NewsFetchCron, true),

		gocron.NewTask(tasks.FetchNewsTask),
	)
	fmt.Printf("JOBS %v", scheduler.Jobs()[0].ID())

	if subscriptionJoberr != nil || newsFetchJobErr != nil {

		panic(fmt.Sprintf("Error in jobs %v", newsFetchJobErr))
	}

	scheduler.Start()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

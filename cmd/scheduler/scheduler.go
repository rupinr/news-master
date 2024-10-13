package main

import (
	"fmt"
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
		gocron.CronJob(
			os.Getenv("SUBSCRIPTION_MAIL_CRON"), false,
		),
		gocron.NewTask(
			tasks.SendNewsletter,
		),
	)

	_, newsFetchJobErr := scheduler.NewJob(
		gocron.CronJob(
			os.Getenv("NEWS_FETCH_CRON"), true,
		),
		gocron.NewTask(tasks.FetchNewsTask),
	)
	if subscriptionJoberr != nil || newsFetchJobErr != nil {
		panic(fmt.Sprintf("Error in jobs %v, %v", subscriptionJoberr, newsFetchJobErr))
	}

	scheduler.Start()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

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
	scheduler, schedulerErr := gocron.NewScheduler()
	if schedulerErr != nil {
		panic(fmt.Sprintf("Unable to start scheduler %v ", schedulerErr.Error()))
	}
	_, subscriptionJoberr := scheduler.NewJob(gocron.CronJob(app.Config.SubscriptionMailCron, false), gocron.NewTask(tasks.SendNewsletter))

	_, newsFetchJobErr := scheduler.NewJob(gocron.CronJob(app.Config.NewsFetchCron, false), gocron.NewTask(tasks.FetchNewsTask))

	_, cleanupJobErr := scheduler.NewJob(gocron.CronJob("0 0 * * *", false), gocron.NewTask(tasks.CleanUp))

	if subscriptionJoberr != nil {
		panic(fmt.Sprintf("Error in subscription job %v ", subscriptionJoberr))
	}

	if newsFetchJobErr != nil {
		panic(fmt.Sprintf("Error in newsFetch job %v ", newsFetchJobErr))
	}

	if cleanupJobErr != nil {
		panic(fmt.Sprintf("Error in cleanup job %v ", cleanupJobErr))
	}

	scheduler.Start()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

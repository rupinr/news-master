package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func main() {
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
			func(a string, b int) {
				fmt.Println("running first  function")
			},
			"hello",
			1,
		),
	)
	if err != nil {
		// handle error
	}
	// each job has a unique id
	fmt.Println(j.ID())

	j2, err2 := s.NewJob(
		gocron.CronJob(
			"* * * * * *", true,
		),
		gocron.NewTask(
			func(a string, b int) {
				fmt.Println("running second function")
			},
			"hello",
			1,
		),
	)
	if err2 != nil {
		// handle error
	}
	// each job has a unique id
	fmt.Println(j2.ID())

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

package main

import (
	"acmilanbot/cli/api"
	"log"
)

func main() {
	err := api.ScheduleJobsSample()
	if err != nil {
		log.Fatal(err)
	}
}

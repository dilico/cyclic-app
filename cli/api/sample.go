package api

import (
	"acmilanbot/api"
	"acmilanbot/cfg"
	"acmilanbot/espn"
	"time"
)

func ScheduleJobsSample() error {
	config, err := cfg.Load()
	if err != nil {
		return err
	}
	config.URL = "https://google.com"
	f := espn.FixtureEvent{
		ID:   "679325",
		Date: time.Now().Add(24 * time.Hour),
		Teams: []espn.Team{
			{
				ID:     "103",
				Name:   "AC Milan",
				IsHome: true,
			},
			{
				ID:     "108",
				Name:   "Udinese",
				IsHome: false,
			},
		},
		Completed: false,
		Venue: espn.Venue{
			Name: "Giuseppe Meazza",
			Address: espn.Address{
				City:    "Milano",
				Country: "Italy",
			},
		},
	}
	return api.ScheduleJobs(f, config)
}

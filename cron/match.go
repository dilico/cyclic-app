package cron

import (
	"acmilanbot/cfg"
	"acmilanbot/espn"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func CreatePreMatchJob(fixture espn.FixtureEvent, config cfg.Configuration) error {
	jobDate := fixture.Date.
		Add(-time.Duration(config.PreMatchHours) * time.Hour)
	jobTitle := fmt.Sprintf(
		"Pre-match: %s - %s",
		espn.GetHomeTeam(fixture.Teams),
		espn.GetAwayTeam(fixture.Teams))
	jobURL := fmt.Sprintf("%s/pre-match", config.URL)
	extendedData, err := json.Marshal(fixture)
	if err != nil {
		return err
	}
	job := newJob(jobTitle, jobURL, []int{jobDate.Hour()}, []int{jobDate.Day()}, []int{jobDate.Minute()}, []int{int(jobDate.Month())}, fixture.Date.Add(24*time.Hour), string(extendedData))
	log.Printf(
		"Creating pre-match job at %s",
		jobDate.Format("2006-01-02T15:04:05"),
	)
	return CreateJob(job, config)
}

func CreateMatchJob(fixture espn.FixtureEvent, config cfg.Configuration) error {
	jobDate := fixture.Date.Add(-30 * time.Minute)
	jobTitle := fmt.Sprintf(
		"Match: %s - %s",
		espn.GetHomeTeam(fixture.Teams),
		espn.GetAwayTeam(fixture.Teams))
	jobURL := fmt.Sprintf("%s/match", config.URL)
	extendedData, err := json.Marshal(fixture)
	if err != nil {
		return err
	}
	job := newJob(jobTitle, jobURL, []int{jobDate.Hour()}, []int{jobDate.Day()}, []int{jobDate.Minute()}, []int{int(jobDate.Month())}, fixture.Date.Add(24*time.Hour), string(extendedData))
	log.Printf(
		"Creating match job at %s",
		jobDate.Format("2006-01-02T15:04:05"),
	)
	return CreateJob(job, config)
}

func CreatePostMatchJob(fixture espn.FixtureEvent, config cfg.Configuration) error {
	maxHours := 3
	getCronHours := func() []int {
		hours := []int{}
		for i := 0; i <= maxHours; i++ {
			d := fixture.Date.Add(time.Duration(i) * time.Hour)
			hours = append(hours, d.Hour())
		}
		return hours
	}
	getCronDays := func() []int {
		days := []int{fixture.Date.Day()}
		d := fixture.Date.Add(time.Duration(maxHours) * time.Hour)
		if !contains(days, d.Day()) {
			days = append(days, d.Day())
		}
		return days
	}
	jobDate := fixture.Date
	jobTitle := fmt.Sprintf(
		"Post-Match: %s - %s",
		espn.GetHomeTeam(fixture.Teams),
		espn.GetAwayTeam(fixture.Teams))
	jobURL := fmt.Sprintf("%s/post-match", config.URL)
	extendedData, err := json.Marshal(fixture)
	if err != nil {
		return err
	}
	job := newJob(jobTitle, jobURL, getCronHours(), getCronDays(), []int{-1}, []int{int(jobDate.Month())}, fixture.Date.Add(24*time.Hour), string(extendedData))
	log.Printf(
		"Creating post-match job at %s",
		jobDate.Format("2006-01-02T15:04:05"),
	)
	return CreateJob(job, config)
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

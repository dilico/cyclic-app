package api

import (
	"acmilanbot/cfg"
	"acmilanbot/cron"
	"acmilanbot/espn"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (s *Server) handleJobsScheduling(w http.ResponseWriter, r *http.Request) {
	fixtures, err := getCandidateFixtures(s.config)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if len(fixtures) == 0 {
		w.WriteHeader(http.StatusNoContent)
		w.Write(nil)
		return
	}
	if len(fixtures) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("too many fixtures"))
		return
	}
	s.config.URL = fmt.Sprintf("%s://%s", r.URL.Scheme, r.Host)
	err = ScheduleJobs(fixtures[0], s.config)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(nil)
}

func getCandidateFixtures(config cfg.Configuration) ([]espn.FixtureEvent, error) {
	candidates := []espn.FixtureEvent{}
	fixtures, err := espn.GetFixtures(config)
	if err != nil {
		return candidates, err
	}
	for _, e := range fixtures.Events {
		if e.Completed {
			continue
		}
		start := time.Now().UTC().
			Add(time.Duration(config.PreMatchHours)*time.Hour + 1*time.Minute)
		end := time.Now().UTC().
			Add(time.Duration(config.PreMatchHours)*time.Hour + 24*time.Hour)
		if e.Date.After(start) && e.Date.Before(end) {
			candidates = append(candidates, e)
		}
	}
	return candidates, nil
}

func ScheduleJobs(fixture espn.FixtureEvent, config cfg.Configuration) error {
	homeTeam := espn.GetHomeTeam(fixture.Teams)
	awayTeam := espn.GetAwayTeam(fixture.Teams)
	log.Printf(
		"New match %s - %s on %s",
		homeTeam,
		awayTeam,
		fixture.Date.Format("2006-01-02T15:04:05"),
	)
	err := cron.DeleteAllJobs(config)
	if err != nil {
		return err
	}
	err = cron.CreatePreMatchJob(fixture, config)
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Minute)
	err = cron.CreateMatchJob(fixture, config)
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Minute)
	return cron.CreatePostMatchJob(fixture, config)
}

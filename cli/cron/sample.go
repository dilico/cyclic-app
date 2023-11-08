package cron

import (
	"acmilanbot/cfg"
	"acmilanbot/cron"
	"log"
)

func GetAllJobsSample() error {
	config, err := cfg.Load()
	if err != nil {
		return err
	}
	jj, err := cron.GetAllJobs(config)
	if err != nil {
		return err
	}
	log.Println(jj)
	return nil
}

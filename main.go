package main

import (
	"acmilanbot/api"
	"acmilanbot/cfg"
	"log"
)

func main() {
	config, err := cfg.Load()
	if err != nil {
		log.Fatal(err)
	}
	s := api.NewServer(config)
	s.Start()
}

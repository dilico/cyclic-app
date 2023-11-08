package cfg

import (
	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Port          int
	PreMatchHours int
	URL           string
	ESPN          ESPN
	Cron          Cron
}

type ESPN struct {
	ACMilanID   string
	FixturesURL string
}

type Cron struct {
	URL    string
	APIKey string
}

func Load() (Configuration, error) {
	c := Configuration{}
	_, err := toml.DecodeFile("config.dev.toml", &c)
	if err != nil {
		_, err = toml.DecodeFile("config.prod.toml", &c)
	}
	return c, err
}

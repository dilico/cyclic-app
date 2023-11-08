package cron

import (
	"acmilanbot/cfg"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

type jobsResponse struct {
	Jobs []Job `json:"jobs"`
}

type Job struct {
	ID             *int                    `json:"jobId,omitempty"`
	Enabled        bool                    `json:"enabled"`
	Title          string                  `json:"title"`
	URL            string                  `json:"url"`
	Schedule       JobSchedule             `json:"schedule,omitempty"`
	RequestMethod  int                     `json:"requestMethod"`
	RequestTimeout int                     `json:"requestTimeout"`
	Notification   JobNotificationSettings `json:"notification,omitempty"`
	ExtendedData   JobExtendedData         `json:"extendedData,omitempty"`
}

type JobSchedule struct {
	Timezone  string `json:"timezone"`
	Hours     []int  `json:"hours"`
	Days      []int  `json:"mdays"`
	Minutes   []int  `json:"minutes"`
	Months    []int  `json:"months"`
	Weekdays  []int  `json:"wdays"`
	ExpiresAt int    `json:"expiresAt"`
}

type JobNotificationSettings struct {
	OnFailure bool `json:"onFailure"`
}

type JobExtendedData struct {
	Body string `json:"body"`
}

type createJobRequest struct {
	Job Job `json:"job"`
}

func newJob(title, url string, hours, days, minutes, months []int, expiresAt time.Time, extendedData string) Job {
	expiration, err := strconv.Atoi(expiresAt.Format("20060102150405"))
	if err != nil {
		expiration = 0
	}
	return Job{
		Enabled: true,
		Title:   title,
		URL:     url,
		Schedule: JobSchedule{
			Timezone:  "UTC",
			Hours:     hours,
			Days:      days,
			Minutes:   minutes,
			Months:    months,
			Weekdays:  []int{-1},
			ExpiresAt: expiration,
		},
		RequestMethod:  1,
		RequestTimeout: 60,
		Notification: JobNotificationSettings{
			OnFailure: true,
		},
		ExtendedData: JobExtendedData{
			Body: extendedData,
		},
	}
}

func GetAllJobs(config cfg.Configuration) ([]Job, error) {
	u, err := url.Parse(config.Cron.URL)
	if err != nil {
		return []Job{}, err
	}
	u.Path = path.Join(u.Path, "jobs")
	request, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return []Job{}, err
	}
	setRequestHeaders(request, config)
	httpClient := http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return []Job{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return []Job{}, fmt.Errorf("cannot get cron jobs [%s]", response.Status)
	}
	jr := jobsResponse{}
	err = json.
		NewDecoder(response.Body).
		Decode(&jr)
	if err != nil {
		return []Job{}, err
	}
	return jr.Jobs, nil
}

func DeleteAllJobs(config cfg.Configuration) error {
	jobs, err := GetAllJobs(config)
	if err != nil {
		return err
	}
	for _, j := range jobs {
		err := DeleteJob(j, config)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteJob(job Job, config cfg.Configuration) error {
	u, err := url.Parse(config.Cron.URL)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, "jobs", strconv.Itoa(*job.ID))
	request, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return err
	}
	setRequestHeaders(request, config)
	httpClient := http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot delete cron job %s [%s]", job.Title, response.Status)
	}
	return nil
}

func CreateJob(job Job, config cfg.Configuration) error {
	u, err := url.Parse(config.Cron.URL)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, "jobs")
	createJobRequest := createJobRequest{
		Job: job,
	}
	requestBody, err := json.Marshal(createJobRequest)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("PUT", u.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	setRequestHeaders(request, config)
	httpClient := http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot create cron job %s [%s]", job.Title, response.Status)
	}
	return nil
}

func setRequestHeaders(r *http.Request, config cfg.Configuration) {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Cron.APIKey))
	r.Header.Set("Content-Type", "application/json")
}

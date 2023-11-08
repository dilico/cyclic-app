package espn

import (
	"acmilanbot/cfg"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type fixturesData struct {
	Page fixturesPage `json:"page"`
}

type fixturesPage struct {
	Content fixturesContent `json:"content"`
}

type fixturesContent struct {
	Fixtures Fixtures `json:"fixtures"`
}

type Fixtures struct {
	Events []FixtureEvent `json:"events"`
}

type FixtureEvent struct {
	ID        string `json:"id"`
	Date      time.Time
	Teams     []Team `json:"teams"`
	Completed bool   `json:"completed"`
	Venue     Venue  `json:"venue"`
}

type Team struct {
	ID     string `json:"id"`
	Name   string `json:"displayName"`
	IsHome bool   `json:"isHome"`
}

type Venue struct {
	Name    string  `json:"fullName"`
	Address Address `json:"address"`
}

type Address struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

func (fe *FixtureEvent) UnmarshalJSON(data []byte) error {
	type Alias FixtureEvent
	aux := &struct {
		Date string `json:"date"`
		*Alias
	}{
		Alias: (*Alias)(fe),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	d, err := time.Parse("2006-01-02T15:04Z", aux.Date)
	if err != nil {
		return err
	}
	fe.Date = d
	return nil
}

func GetFixtures(config cfg.Configuration) (Fixtures, error) {
	u, err := url.Parse(config.ESPN.FixturesURL)
	if err != nil {
		return Fixtures{}, err
	}
	u.Path = path.Join(u.Path, config.ESPN.ACMilanID)
	request, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return Fixtures{}, err
	}
	httpClient := http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return Fixtures{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return Fixtures{}, fmt.Errorf("cannot get espn fixtures page [%s]", response.Status)
	}
	doc, err := html.Parse(response.Body)
	if err != nil {
		return Fixtures{}, err
	}
	return getFixturesFromHTML(doc)
}

func getFixturesFromHTML(doc *html.Node) (Fixtures, error) {
	fixturesJSON := getFixturesJSONFromHTML(doc)
	if len(fixturesJSON) == 0 {
		return Fixtures{}, errors.New("cannot find espn fixtures")
	}
	fd := fixturesData{}
	err := json.
		NewDecoder(strings.NewReader(fixturesJSON)).
		Decode(&fd)
	if err != nil {
		return Fixtures{}, err
	}
	return fd.Page.Content.Fixtures, nil
}

func getFixturesJSONFromHTML(doc *html.Node) string {
	var fixturesJSON string
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if len(fixturesJSON) > 0 {
			return
		}
		fixturesJSON = extractFixturesJSON(node)
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	return fixturesJSON
}

func extractFixturesJSON(node *html.Node) string {
	if node.Type == html.ElementNode && node.Data == "script" {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode && strings.Contains(child.Data, "window['__espnfitt__']") {
				return extractJSONFromNode(child)
			}
		}
	}
	return ""
}

func extractJSONFromNode(node *html.Node) string {
	rePattern := `.*window\['__espnfitt__']=(.*)`
	re := regexp.MustCompile(rePattern)
	matches := re.FindStringSubmatch(node.Data)
	if len(matches) < 1 {
		return ""
	}
	return strings.TrimSuffix(matches[1], ";")
}

func GetHomeTeam(teams []Team) string {
	for _, t := range teams {
		if t.IsHome {
			return t.Name
		}
	}
	return ""
}

func GetAwayTeam(teams []Team) string {
	for _, t := range teams {
		if !t.IsHome {
			return t.Name
		}
	}
	return ""
}

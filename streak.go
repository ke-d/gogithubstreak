package streak

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Streak with count and the dates of the streak
type Streak struct {
	From  time.Time `json:"from"`
	To    time.Time `json:"to"`
	Count int       `json:"count"`
}

// Client with http Client for dependency injection
type Client struct {
	Client  *http.Client
	BaseURL string
}

func getStreakFromCalendar(doc *goquery.Document) (Streak, Streak, error) {

	curStreak := Streak{Count: 0}
	longestStreak := Streak{Count: 0}
	isStreak := false

	allDays := doc.Find(".day")

	allDays.Each(func(i int, days *goquery.Selection) {
		count, existsCount := days.Attr("data-count")
		date, existsDate := days.Attr("data-date")

		if !existsCount || !existsDate {
			return
		}

		parsedDate, err := time.Parse("2006-1-2", date)

		if count != "0" && !isStreak {
			// New streak
			curStreak.Count++

			if err == nil {
				curStreak.From = parsedDate
			}
			isStreak = true
		} else if count != "0" && isStreak {
			// Still has a streak
			curStreak.Count++
		} else if count == "0" && isStreak && allDays.Length()-1 != i {
			// Lost streak
			if err == nil {
				curStreak.To = parsedDate.AddDate(0, 0, -1)
			}
			if longestStreak.Count < curStreak.Count {
				longestStreak = curStreak
			}
			curStreak = Streak{Count: 0}
			isStreak = false
		}

	})

	// Check last day
	last := allDays.Last()
	count, existsCount := last.Attr("data-count")
	lastDate, existsDate := last.Attr("data-date")
	if existsCount && existsDate {
		parsedLastDate, err := time.Parse("2006-1-2", lastDate)
		if isStreak && err == nil && count != "0" {
			// If the last day has a streak
			curStreak.To = parsedLastDate

		} else if isStreak && err == nil && count == "0" {
			// If the user hasn't commited on the last day, set current streak to the day before,
			// since it doesn't mean that the user lost that streak if the user didn't commit yet on the current day
			curStreak.To = parsedLastDate.AddDate(0, 0, -1)
		}
	}

	// If the current streak is the longest streak, then return that for longest
	if curStreak.Count > longestStreak.Count {
		return curStreak, curStreak, nil
	}
	return curStreak, longestStreak, nil
}

func getCalendarFromGitHub(client *Client, username string, date time.Time) (*http.Response, error) {
	resp, err := client.Client.Get(client.BaseURL + "/users/" + username + "/contributions?to=" + date.Format("2006-1-2"))
	if resp.StatusCode != 200 {
		return resp, errors.New("Cannot get calendar")
	}
	return resp, err
}

func getContributions(doc *goquery.Document) int {
	reg, _ := regexp.Compile(`([\d]*) contributions`)
	matchArr := reg.FindStringSubmatch(doc.Find(".f4").Text())

	if len(matchArr) < 1 {
		return 0
	}
	numOfContributionsStr := matchArr[1]
	// Reg exp match should only have numbers at this point
	numOfContributions, _ := strconv.Atoi(numOfContributionsStr)
	return numOfContributions
}

// FindStreakInPastYear returns the Current streak as the first return and the Longest streak in the second return as well as a potential error.
func FindStreakInPastYear(client *Client, username string) (Streak, Streak, error) {
	now := time.Now()
	resp, err := getCalendarFromGitHub(client, username, now)
	if err != nil {
		return Streak{}, Streak{}, errors.New("Cannot get calendar")
	}

	defer resp.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return Streak{}, Streak{}, errors.New("Could not load calendar")
	}

	numOfContributions := getContributions(doc)

	if numOfContributions != 0 {
		return getStreakFromCalendar(doc)
	}
	return Streak{}, Streak{}, errors.New("No contributions")
}

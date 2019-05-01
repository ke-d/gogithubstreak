package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Streak - The Streak from Date to To Date
type Streak struct {
	From  time.Time `json:"from"`
	To    time.Time `json:"to"`
	Count int       `json:"count"`
}

func getStreakFromCalendar(doc *goquery.Document) (Streak, Streak, error) {

	curStreak := Streak{Count: 0}
	longestStreak := Streak{Count: 0}
	isStreak := false

	allDays := doc.Find(".day")

	allDays.Each(func(i int, s *goquery.Selection) {
		count, exists1 := s.Attr("data-count")
		date, exists2 := s.Attr("data-date")

		parsedDate, err := time.Parse("2006-1-2", date)

		if !exists1 || !exists2 {
			return
		}
		// fmt.Println(count, isStreak, curStreak, longestStreak)
		if count != "0" && !isStreak {
			curStreak.Count++

			if err == nil {
				curStreak.From = parsedDate
			}
			isStreak = true
		} else if count != "0" && isStreak {
			curStreak.Count++
		} else if count == "0" && isStreak {
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

	lastDate, exists := allDays.Last().Attr("data-date")
	parsedDate, err := time.Parse("2006-1-2", lastDate)
	if isStreak && exists && err == nil {
		curStreak.To = parsedDate
	}

	if curStreak.Count > longestStreak.Count {
		return curStreak, curStreak, nil
	}
	return curStreak, longestStreak, nil
	// fmt.Fprintf(w, "hello")
}

// https://github.com/users/mrdokenny/contributions?to=2016-1-1

func getCalendarFromGitHub(username string, date time.Time) (*http.Response, error) {
	fmt.Println(username)
	resp, err := http.Get("https://github.com/users/" + username + "/contributions?to=" + date.Format("2006-1-2"))

	return resp, err
}

func findStreak(username string) (Streak, Streak, error) {
	now := time.Now()
	resp, err := getCalendarFromGitHub(username, now)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Printf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	// Load the HTML document
	doc, err2 := goquery.NewDocumentFromReader(resp.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	reg, _ := regexp.Compile(`([\d]*) contributions`)
	numOfContributions := reg.FindStringSubmatch(doc.Find(".f4").Text())[1]

	if numOfContributions != "0" {
		return getStreakFromCalendar(doc)
	}
	return Streak{}, Streak{}, errors.New("No contributions")
}

func handler(w http.ResponseWriter, r *http.Request) {
	uriSegments := strings.Split(r.URL.Path, "/")
	username := strings.ToLower(uriSegments[1])
	if username == "" {
		http.NotFound(w, r)
		return
	}
	curStreak, longestStreak, _ := findStreak(username)
	fmt.Println(curStreak, longestStreak)

}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
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

func getStreakFromCalendar(resp *http.Response) {

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Printf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	// Load the HTML document
	doc, err2 := goquery.NewDocumentFromReader(resp.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	r, _ := regexp.Compile(`([\d]*) contributions`)
	numOfContributions := r.FindStringSubmatch(doc.Find(".f4").Text())[1]
	fmt.Println("dasd", numOfContributions)
	if numOfContributions == "0" {
		return
	}

	curStreak := 0
	longestStreak := 0
	isStreak := false

	allDays := doc.Find(".day")

	firstDate, exists := allDays.First().Attr("data-date")

	if exists {
		fmt.Println(firstDate)
	}

	allDays.Each(func(i int, s *goquery.Selection) {
		count, exists := s.Attr("data-count")
		if !exists {
			return
		}
		fmt.Println(count, isStreak)
		if count != "0" && !isStreak {
			curStreak++
			isStreak = true
		} else if count != "0" && isStreak {
			curStreak++
		} else if count == "0" && isStreak {
			longestStreak = curStreak
			curStreak = 0
			isStreak = false
		}

	})
	fmt.Println("cur", curStreak)
	fmt.Println("longest", longestStreak)
	// fmt.Fprintf(w, "hello")
}

// https://github.com/users/mrdokenny/contributions?to=2016-1-1

func getStreak(username string, date time.Time) (*http.Response, error) {
	fmt.Println(username)
	resp, err := http.Get("https://github.com/users/" + username + "/contributions?to=" + date.Format("2006-1-2"))

	return resp, err
}

func handler(w http.ResponseWriter, r *http.Request) {
	uriSegments := strings.Split(r.URL.Path, "/")
	username := strings.ToLower(uriSegments[1])
	if username == "" {
		http.NotFound(w, r)
		return
	}

	now := time.Now()
	resp, err := getStreak(username, now)
	if err != nil {
		fmt.Println(err)
	}
	getStreakFromCalendar(resp)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

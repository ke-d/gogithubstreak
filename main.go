package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func handler(w http.ResponseWriter, r *http.Request) {
	uriSegments := strings.Split(r.URL.Path, "/")
	username := strings.ToLower(uriSegments[1])
	if username == "" {
		http.NotFound(w, r)
		return
	}
	fmt.Println(username)
	resp, err1 := http.Get("https://github.com/users/" + username + "/contributions")

	if err1 != nil {
		http.NotFound(w, r)
		fmt.Printf("%s", err1)
		return
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

	curStreak := 0
	longestStreak := 0
	isStreak := false
	doc.Find(".day").Each(func(i int, s *goquery.Selection) {
		count, exists := s.Attr("data-count")
		if !exists {
			return
		}
		num, err := strconv.Atoi(count)
		if err != nil {
			// handle error
			fmt.Println(err)
			// os.Exit(2)
		}
		if num != 0 && !isStreak {
			curStreak++
			isStreak = true
		} else if num != 0 && isStreak {
			curStreak++
		} else if num == 0 && isStreak {
			longestStreak = curStreak
			curStreak = 0
			isStreak = false
		}

	})
	fmt.Println(curStreak)
	fmt.Println(longestStreak)
	fmt.Fprintf(w, "hello")
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

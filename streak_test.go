package streak

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func Test_getStreakFromCalendar(t *testing.T) {
	// doc1 has a day without data-count and data-date which should be ignored
	doc1, _ := goquery.NewDocumentFromReader(strings.NewReader(`
	<g transform="translate(0, 0)">
	<rect class="day" width="8" height="8" x="11" y="30" fill="#ebedf0"></rect>
	<rect class="day" width="8" height="8" x="11" y="30" fill="#ebedf0" data-count="1" data-date="2018-05-02"></rect>
	<rect class="day" width="8" height="8" x="11" y="40" fill="#c6e48b" data-count="3" data-date="2018-05-03"></rect>
	<rect class="day" width="8" height="8" x="11" y="50" fill="#c6e48b" data-count="1" data-date="2018-05-04"></rect>
	<rect class="day" width="8" height="8" x="11" y="60" fill="#c6e48b" data-count="1" data-date="2018-05-05"></rect>
	<rect class="day" width="8" height="8" x="10" y="0" fill="#ebedf0" data-count="0" data-date="2018-05-06"></rect>
	<rect class="day" width="8" height="8" x="10" y="10" fill="#ebedf0" data-count="0" data-date="2018-05-07"></rect>
	<rect class="day" width="8" height="8" x="10" y="20" fill="#ebedf0" data-count="0" data-date="2018-05-08"></rect>
	<rect class="day" width="8" height="8" x="10" y="30" fill="#c6e48b" data-count="1" data-date="2018-05-09"></rect>
	<rect class="day" width="8" height="8" x="10" y="40" fill="#ebedf0" data-count="0" data-date="2018-05-10"></rect>
	<rect class="day" width="8" height="8" x="10" y="50" fill="#239a3b" data-count="7" data-date="2018-05-11"></rect>
	<rect class="day" width="8" height="8" x="10" y="60" fill="#7bc96f" data-count="5" data-date="2018-05-12"></rect>
</g>
`))
	doc2, _ := goquery.NewDocumentFromReader(strings.NewReader(`
	<g transform="translate(0, 0)">
	<rect class="day" width="8" height="8" x="11" y="30" fill="#ebedf0" data-count="1" data-date="2018-05-02"></rect>
	<rect class="day" width="8" height="8" x="11" y="40" fill="#c6e48b" data-count="1" data-date="2018-05-03"></rect>
	<rect class="day" width="8" height="8" x="11" y="50" fill="#c6e48b" data-count="0" data-date="2018-05-04"></rect>
	<rect class="day" width="8" height="8" x="11" y="60" fill="#c6e48b" data-count="0" data-date="2018-05-05"></rect>
	<rect class="day" width="8" height="8" x="10" y="0" fill="#ebedf0" data-count="0" data-date="2018-05-06"></rect>
	<rect class="day" width="8" height="8" x="10" y="10" fill="#ebedf0" data-count="0" data-date="2018-05-07"></rect>
	<rect class="day" width="8" height="8" x="10" y="20" fill="#ebedf0" data-count="0" data-date="2018-05-08"></rect>
	<rect class="day" width="8" height="8" x="10" y="30" fill="#c6e48b" data-count="1" data-date="2018-05-09"></rect>
	<rect class="day" width="8" height="8" x="10" y="40" fill="#ebedf0" data-count="2" data-date="2018-05-10"></rect>
	<rect class="day" width="8" height="8" x="10" y="50" fill="#239a3b" data-count="7" data-date="2018-05-11"></rect>
	<rect class="day" width="8" height="8" x="10" y="60" fill="#7bc96f" data-count="5" data-date="2018-05-12"></rect>
	<rect class="day" width="8" height="8" x="10" y="60" fill="#7bc96f" data-count="2" data-date="2018-05-13"></rect>
</g>
`))
	doc3, _ := goquery.NewDocumentFromReader(strings.NewReader(`
	<g transform="translate(0, 0)">
	<rect class="day" width="8" height="8" x="11" y="30" fill="#ebedf0" data-count="1" data-date="2018-05-02"></rect>
	<rect class="day" width="8" height="8" x="11" y="40" fill="#c6e48b" data-count="1" data-date="2018-05-03"></rect>
	<rect class="day" width="8" height="8" x="11" y="50" fill="#c6e48b" data-count="0" data-date="2018-05-04"></rect>
	<rect class="day" width="8" height="8" x="11" y="60" fill="#c6e48b" data-count="0" data-date="2018-05-05"></rect>
	<rect class="day" width="8" height="8" x="10" y="0" fill="#ebedf0" data-count="0" data-date="2018-05-06"></rect>
	<rect class="day" width="8" height="8" x="10" y="10" fill="#ebedf0" data-count="0" data-date="2018-05-07"></rect>
	<rect class="day" width="8" height="8" x="10" y="20" fill="#ebedf0" data-count="0" data-date="2018-05-08"></rect>
	<rect class="day" width="8" height="8" x="10" y="30" fill="#c6e48b" data-count="1" data-date="2018-05-09"></rect>
	<rect class="day" width="8" height="8" x="10" y="40" fill="#ebedf0" data-count="2" data-date="2018-05-10"></rect>
	<rect class="day" width="8" height="8" x="10" y="50" fill="#239a3b" data-count="7" data-date="2018-05-11"></rect>
	<rect class="day" width="8" height="8" x="10" y="60" fill="#7bc96f" data-count="5" data-date="2018-05-12"></rect>
	<rect class="day" width="8" height="8" x="10" y="60" fill="#7bc96f" data-count="2" data-date="2018-05-13"></rect>
	<rect class="day" width="8" height="8" x="10" y="60" fill="#7bc96f" data-count="0" data-date="2018-05-14"></rect>
</g>
`))
	doc4, _ := goquery.NewDocumentFromReader(strings.NewReader(`
	<g transform="translate(0, 0)">
	<rect class="day" width="8" height="8" x="11" y="30" fill="#ebedf0" data-count="1" data-date="2018-05-02"></rect>
	<rect class="day" width="8" height="8" x="11" y="40" fill="#c6e48b" data-count="1" data-date="2018-05-03"></rect>
	<rect class="day" width="8" height="8" x="11" y="50" fill="#c6e48b" data-count="0" data-date="2018-05-04"></rect>
	<rect class="day" width="8" height="8" x="11" y="60" fill="#c6e48b" data-count="0" data-date="2018-05-05"></rect>
	<rect class="day" width="8" height="8" x="10" y="0" fill="#ebedf0" data-count="0" data-date="2018-05-06"></rect>
	<rect class="day" width="8" height="8" x="10" y="10" fill="#ebedf0" data-count="0" data-date="2018-05-07"></rect>
	<rect class="day" width="8" height="8" x="10" y="20" fill="#ebedf0" data-count="0" data-date="2018-05-08"></rect>
	<rect class="day" width="8" height="8" x="10" y="30" fill="#c6e48b" data-count="1" data-date="2018-05-09"></rect>
	<rect class="day" width="8" height="8" x="10" y="40" fill="#ebedf0" data-count="2" data-date="2018-05-10"></rect>
	<rect class="day" width="8" height="8" x="10" y="50" fill="#239a3b" data-count="7" data-date="2018-05-11"></rect>
	<rect class="day" width="8" height="8" x="10" y="60" fill="#7bc96f" data-count="5" data-date="2018-05-12"></rect>
	<rect class="day" width="8" height="8" x="10" y="60" fill="#7bc96f" data-count="2" data-date="2018-05-13"></rect>
	<rect class="day" width="8" height="8" x="10" y="60" fill="#7bc96f" data-count="0" data-date="2018-05-14"></rect>
	<rect class="day" width="8" height="8" x="10" y="60" fill="#7bc96f" data-count="0" data-date="2018-05-15"></rect>
</g>
`))

	type args struct {
		doc *goquery.Document
	}
	tests := []struct {
		name    string
		args    args
		want    Streak
		want1   Streak
		wantErr bool
	}{
		{"Basic Test", args{doc: doc1},
			Streak{
				Count: 2,
				From:  time.Date(2018, 5, 11, 0, 0, 0, 0, time.UTC),
				To:    time.Date(2018, 5, 12, 0, 0, 0, 0, time.UTC),
			},
			Streak{
				Count: 4,
				From:  time.Date(2018, 5, 2, 0, 0, 0, 0, time.UTC),
				To:    time.Date(2018, 5, 5, 0, 0, 0, 0, time.UTC),
			},
			false},
		{"Test for keeping current streak", args{doc: doc2},
			Streak{
				Count: 5,
				From:  time.Date(2018, 5, 9, 0, 0, 0, 0, time.UTC),
				To:    time.Date(2018, 5, 13, 0, 0, 0, 0, time.UTC),
			},
			Streak{
				Count: 5,
				From:  time.Date(2018, 5, 9, 0, 0, 0, 0, time.UTC),
				To:    time.Date(2018, 5, 13, 0, 0, 0, 0, time.UTC),
			},
			false},
		{"Test for missing one day of current streak", args{doc: doc3},
			Streak{
				Count: 5,
				From:  time.Date(2018, 5, 9, 0, 0, 0, 0, time.UTC),
				To:    time.Date(2018, 5, 13, 0, 0, 0, 0, time.UTC),
			},
			Streak{
				Count: 5,
				From:  time.Date(2018, 5, 9, 0, 0, 0, 0, time.UTC),
				To:    time.Date(2018, 5, 13, 0, 0, 0, 0, time.UTC),
			},
			false},
		{"Test for missing two days of current streak", args{doc: doc4},
			Streak{
				Count: 0,
				From:  time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
				To:    time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			Streak{
				Count: 5,
				From:  time.Date(2018, 5, 9, 0, 0, 0, 0, time.UTC),
				To:    time.Date(2018, 5, 13, 0, 0, 0, 0, time.UTC),
			},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getStreakFromCalendar(tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("getStreakFromCalendar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getStreakFromCalendar() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getStreakFromCalendar() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getContributions(t *testing.T) {
	doc1, _ := goquery.NewDocumentFromReader(strings.NewReader(`
		<div>
			Test
		</div>
	`))
	doc2, _ := goquery.NewDocumentFromReader(strings.NewReader(`
	<h2 class="f4 text-normal mb-2">
		520 contributions
			in the last year
	</h2>
	`))
	doc3, _ := goquery.NewDocumentFromReader(strings.NewReader(`
	<h2 class="f4 text-normal mb-2">
		5-20 contributions
			in the last year
	</h2>
	`))
	type args struct {
		doc *goquery.Document
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Non contribution page", args{doc: doc1}, 0},
		{"Page with contributions", args{doc: doc2}, 520},
		{"Page with malformed contributions", args{doc: doc3}, 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getContributions(tt.args.doc); got != tt.want {
				t.Errorf("getContributions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCalendarFromGitHub(t *testing.T) {
	tsOkResponse := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(`OK`))
	}))
	defer tsOkResponse.Close()
	tsBadResponse := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`Not Found`))
	}))
	defer tsBadResponse.Close()
	type args struct {
		client   Client
		username string
		date     time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{"Ok Response", args{client: Client{tsOkResponse.Client(), tsOkResponse.URL}, username: "any", date: time.Now()},
			&http.Response{
				StatusCode: 200},
			false},
		{"Bad Response", args{client: Client{tsBadResponse.Client(), tsBadResponse.URL}, username: "any", date: time.Now()},
			&http.Response{
				StatusCode: 404},
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCalendarFromGitHub(&tt.args.client, tt.args.username, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCalendarFromGitHub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.StatusCode != tt.want.StatusCode {
				t.Errorf("getCalendarFromGitHub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindStreakInPastYear(t *testing.T) {
	tsOkResponse := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(`
		<div class="f4">2 contributions</div>
		<rect class="day" width="8" height="8" x="11" y="30" fill="#ebedf0" data-count="1" data-date="2018-05-02"></rect>
		<rect class="day" width="8" height="8" x="11" y="40" fill="#c6e48b" data-count="1" data-date="2018-05-03"></rect>
		<rect class="day" width="8" height="8" x="11" y="50" fill="#c6e48b" data-count="0" data-date="2018-05-04"></rect>
		`))
	}))
	defer tsOkResponse.Close()

	tsBadResponse := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)

		w.Write([]byte(`Not Found`))
	}))
	defer tsBadResponse.Close()

	tsInvalidHTML := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(`invalid HTML`))
	}))
	defer tsInvalidHTML.Close()

	tsOkNoContributionsResponse := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(`
		<div class="f4">0 contributions</div>
		`))
	}))
	defer tsOkNoContributionsResponse.Close()
	type args struct {
		client   Client
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    Streak
		want1   Streak
		wantErr bool
	}{
		{
			"Test Ok Response with valid HTML",
			args{
				client:   Client{tsOkResponse.Client(), tsOkResponse.URL},
				username: "any",
			},
			Streak{
				Count: 2,
				From:  time.Date(2018, 5, 2, 0, 0, 0, 0, time.UTC),
				To:    time.Date(2018, 5, 3, 0, 0, 0, 0, time.UTC),
			},
			Streak{
				Count: 2,
				From:  time.Date(2018, 5, 2, 0, 0, 0, 0, time.UTC),
				To:    time.Date(2018, 5, 3, 0, 0, 0, 0, time.UTC),
			},
			false,
		},
		{
			"Test not found",
			args{
				client:   Client{tsBadResponse.Client(), tsBadResponse.URL},
				username: "any",
			},
			Streak{},
			Streak{},
			true,
		},
		{
			// Suppose to be testing the goquery.NewDocumentFromReade but it seems to go through
			"Test invalid HTML",
			args{
				client:   Client{tsInvalidHTML.Client(), tsInvalidHTML.URL},
				username: "any",
			},
			Streak{},
			Streak{},
			true,
		},
		{
			"Test no contributions",
			args{
				client:   Client{tsOkNoContributionsResponse.Client(), tsOkNoContributionsResponse.URL},
				username: "any",
			},
			Streak{},
			Streak{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := FindStreakInPastYear(&tt.args.client, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindStreakInPastYear() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindStreakInPastYear() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("FindStreakInPastYear() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

package livehelp

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/udacity/mc/cc"
	"github.com/udacity/mc/csv"
	"github.com/udacity/mc/httpclient"
	"github.com/udacity/mc/mentor"
	"github.com/udacity/mc/payments"
)

// Response represents the response object received from LiveHelp API.
type Response struct {
	Data   Data    `json:"data"`
	Errors []Error `json:"errors"`
}

// Data represents a collection of data provided by the Response.
type Data struct {
	ResolvedSessions []ResolvedSession `json:"resolved_sessions"`
}

// Error represents a collection of errors provided by the Response.
type Error struct {
	Message string `json:"message"`
}

// ResolvedSession represents a LiveHelp resolved session.
type ResolvedSession struct {
	UIDExpert     string    `json:"uid_expert"`
	UIDStudent    string    `json:"uid_student"`
	KeyNanodegree string    `json:"key_nanodegree"`
	DateResolved  time.Time `json:"date_resolved"`
}

// FetchResolvedSessions fetches a resolved session from LiveHelp API using a UID.
func FetchResolvedSessions(isStaging bool, startDate, endDate string) ([]ResolvedSession, error) {
	client, err := httpclient.New(url(isStaging))
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("sessions/resolved?start_date=%s&end_date=%s", startDate, endDate)

	res := Response{}
	err = client.Call(http.MethodGet, endpoint, nil, &res)
	if err != nil {
		return nil, err
	}

	return res.Data.ResolvedSessions, nil
}

// CalculatePayment information for livehelpers during a specified period.
func CalculatePayment(isStaging bool, startDate, endDate string) (string, error) {
	resolvedSessions, err := FetchResolvedSessions(isStaging, startDate, endDate)
	if err != nil {
		return "", err
	}

	// Make set of livehelper UIDs.
	livehelpers := make(map[string]bool)
	for _, session := range resolvedSessions {
		livehelpers[session.UIDExpert] = true
	}

	// Map UID to resolved-session counts by Nanodegree.
	sessionCounter := make(map[string]map[string]int)
	for _, session := range resolvedSessions {
		if sessionCounter[session.UIDExpert] == nil {
			sessionCounter[session.UIDExpert] = make(map[string]int)
		}
		sessionCounter[session.UIDExpert][session.KeyNanodegree]++
	}

	// Fetch mentors.
	mentors, err := mentor.FetchMentors(isStaging)
	if err != nil {
		return "", err
	}

	// Map UIDs to mentors.
	mentorMap := make(map[string]mentor.Mentor)
	for _, mentor := range mentors {
		mentorMap[mentor.UID] = mentor
	}

	// Map countries to payments.
	paymentsInfo := make(map[string]*payments.Data)
	for _, m := range mentors {
		// Only fetch payments info for countries
		// not already fetched.
		if m.Country != "" && paymentsInfo[m.Country] == nil {
			paymentData, err := payments.Fetch(isStaging, m.Country)
			// If err, completely bail from entire payments process.
			if err != nil {
				return "", err
			}
			paymentsInfo[m.Country] = &paymentData
		}
	}

	// Create rows.
	rowsChan := make(chan []string)
	var wg sync.WaitGroup
	for uid := range livehelpers {
		wg.Add(1)
		go createLiveHelperRow(mentorMap[uid], sessionCounter[uid], paymentsInfo, &wg, rowsChan)
	}
	ticker := time.NewTicker(time.Second * 1)
	go func() {
		time.Sleep(time.Second * 1)
		fmt.Print("Processing")
		for range ticker.C {
			fmt.Print(".")
		}
	}()
	go func() {
		wg.Wait()
		close(rowsChan)
		ticker.Stop()
		fmt.Println()
	}()

	var rows [][]string
	for row := range rowsChan {
		rows = append(rows, row)
	}

	// Create CSV and check for errors.
	header := []string{"uid", "first_name", "last_name", "email", "paypal_email", "country", "nanodegree_key", "nanodegree_level", "resolved_sessions", "amount"}
	suffix := fmt.Sprintf("%s_%s", startDate, endDate)
	return csv.CreateCSVFile(liveHelpPaymentsFilename, suffix, header, rows)
}

func createLiveHelperRow(m mentor.Mentor, sessionCounts map[string]int, paymentsInfo map[string]*payments.Data, wg *sync.WaitGroup, rowsChan chan<- []string) {
	defer wg.Done()
	// Fetch user. Skip if user information is not available.
	user, err := cc.FetchUser(m.UID)
	if err != nil {
		log.Fatalf("error: fetching user %s from Classroom %s\n", m.UID, err)
	}

	// Create some rows.
	for ndKey, sessionCount := range sessionCounts {
		ndLevel := paymentsInfo[m.Country].GetNDLevel(ndKey)

		rowsChan <- []string{
			user.ID,
			user.FirstName,
			user.LastName,
			user.Email,
			m.PayPalEmail,
			m.Country,
			ndKey,
			ndLevel,
			fmt.Sprintf("%d", sessionCount),
			displayAmount(m.Country, ndLevel, sessionCount, paymentsInfo[m.Country]),
		}
	}
}

func calculateAmount(ndLevel string, sessionCount int, livehelpSessions *payments.LiveHelpSessions) float32 {
	act := livehelpSessions.ActivityLevels["active"]
	ver := livehelpSessions.ActivityLevels["very-active"]
	ext := livehelpSessions.ActivityLevels["extremely-active"]

	amounts := livehelpSessions.Amounts[ndLevel]

	var amount float32
	for i := 1; i <= sessionCount; i++ {
		switch {
		case i >= ext:
			amount += 1 * amounts["extremely-active"]
		case i >= ver:
			amount += 1 * amounts["very-active"]
		case i >= act:
			amount += 1 * amounts["active"]
		}
	}
	return amount
}

func displayAmount(country, ndLevel string, sessionCount int, paymentsData *payments.Data) string {
	if country == "" {
		return "COUNTRY IS NULL"
	}
	if paymentsData == nil || len(paymentsData.LiveHelpSessions.Amounts) == 0 || len(paymentsData.LiveHelpSessions.Amounts["beginner"]) == 0 {
		return fmt.Sprintf("COUNTRY '%s' IS NOT KNOWN", country)
	}
	return fmt.Sprintf("%.2f", calculateAmount(ndLevel, sessionCount, paymentsData.LiveHelpSessions))
}

func url(isStaging bool) string {
	if isStaging {
		return stagingURL
	}
	return productionURL
}

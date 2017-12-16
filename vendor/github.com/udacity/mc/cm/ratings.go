package cm

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/udacity/mc/cc"
	"github.com/udacity/mc/csv"
	"github.com/udacity/mc/mentor"
	"github.com/udacity/mc/payments"
)

// RatingPayments calculates payment information for classroom mentors during
// a specified period based on ratings.
func RatingPayments(startDate, endDate string, isStaging bool) (string, error) {
	err := checkFileExists(guruRatingsFilename, startDate, endDate)
	if err != nil {
		return "", err
	}
	bonusCountRows, err := readCSV(guruRatingsFilename)
	if err != nil {
		return "", err
	}

	// Bring in user info.
	userRows, err := readCSV(usersFilename)
	if err != nil {
		return "", err
	}

	// Map UIDs to users.
	userMap := make(map[string]cc.User)
	for _, user := range userRows {
		uid, firstName, lastName, email := user[0], user[1], user[2], user[3]
		userMap[uid] = cc.User{
			ID:        uid,
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
		}
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

	// Start ticker.
	ticker := time.NewTicker(time.Second * 1)
	go func() {
		time.Sleep(time.Second * 1)
		fmt.Print("Processing")
		for range ticker.C {
			fmt.Print(".")
		}
	}()

	// Create rows.
	rowsChan := make(chan []string)
	var wg sync.WaitGroup
	for _, bonusCountRow := range bonusCountRows {
		// Skip headers.
		if bonusCountRow[3] == "count" {
			continue
		}
		// Extract and store UID.
		uid := bonusCountRow[0]

		// Only process those records where the UID is in mentorMap.
		if mentorMap[uid].UID != "" {
			wg.Add(1)
			go createClassMentorRatingRow(mentorMap[uid], userMap[uid], bonusCountRow, paymentsInfo, &wg, rowsChan)
		} else {
			fmt.Printf("error: mentor %s info was not available\n", uid)
		}
	}

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
	header := []string{"uid", "first_name", "last_name", "email", "paypal_email", "country", "nanodegree_key", "week", "bonus_count", "amount"}
	suffix := fmt.Sprintf("%s_%s", startDate, endDate)
	return csv.CreateCSVFile(classMentorRatingsPaymentsFilename, suffix, header, rows)
}

func createClassMentorRatingRow(mentor mentor.Mentor, user cc.User, bonusCountRow []string, paymentsInfo map[string]*payments.Data, wg *sync.WaitGroup, rowsChan chan<- []string) {
	defer wg.Done()

	// Extract and store values.
	ndKey, week := bonusCountRow[1], bonusCountRow[2]
	bonusCount, err := strconv.Atoi(bonusCountRow[3])
	if err != nil {
		return
	}

	// Create row.
	rowsChan <- []string{
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		mentor.PayPalEmail,
		mentor.Country,
		ndKey,
		week,
		strconv.Itoa(bonusCount),
		displayRatingAmount(mentor.Country, bonusCount, paymentsInfo[mentor.Country]),
	}
}

func calculateRatingAmount(bonusCount int, classMentorRatings *payments.ClassroomMentorshipRatings) float32 {
	return float32(bonusCount) * classMentorRatings.Amounts["highest"]
}

func displayRatingAmount(country string, bonusCount int, paymentsData *payments.Data) string {
	if country == "" {
		return "COUNTRY IS NULL"
	}
	if paymentsData == nil || len(paymentsData.ClassroomMentorshipRatings.Amounts) == 0 {
		return fmt.Sprintf("COUNTRY '%s' IS NOT KNOWN", country)
	}
	return fmt.Sprintf("%.2f", calculateRatingAmount(bonusCount, paymentsData.ClassroomMentorshipRatings))
}

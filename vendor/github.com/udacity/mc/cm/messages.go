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

// MessagePayments calculates payment information for classroom mentors during
// a specified period based on messages.
func MessagePayments(startDate, endDate string, isStaging bool) (string, error) {
	msgFilenames := []string{
		androidMessagesFilename,
		iOSMessagesFilename,
		webMessagesFilename,
		guruCheckinsFilename,
	}

	// Read values from CSVs.
	var msgCountRows [][]string
	for _, msgFilename := range msgFilenames {
		err := checkFileExists(msgFilename, startDate, endDate)
		if err != nil {
			return "", err
		}
		rows, err := readCSV(msgFilename)
		if err != nil {
			return "", err
		}
		msgCountRows = append(msgCountRows, rows...)
	}

	// Aggregate counts.
	msgCounts := make(map[string]map[string]map[string]map[string]int)
	for _, row := range msgCountRows {
		// Skip headers.
		if row[4] == "count" {
			continue
		}
		// Extract and store values.
		uid, ndKey, studUID, week := row[0], row[1], row[2], row[3]
		count, err := strconv.Atoi(row[4])
		if err != nil {
			return "", fmt.Errorf("%s value: '%s'", err, row[4])
		}
		if msgCounts[uid] == nil {
			msgCounts[uid] = make(map[string]map[string]map[string]int)
		}
		if msgCounts[uid][ndKey] == nil {
			msgCounts[uid][ndKey] = make(map[string]map[string]int)
		}
		if msgCounts[uid][ndKey][studUID] == nil {
			msgCounts[uid][ndKey][studUID] = make(map[string]int)
		}
		msgCounts[uid][ndKey][studUID][week] += count
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
	for uid := range msgCounts {
		// Message records include those sent from mentors to students.
		// Only process those records where the recipient UID is in mentorMap.
		if mentorMap[uid].UID != "" {
			wg.Add(1)
			go createClassMentorMessageRow(mentorMap[uid], userMap[uid], msgCounts[uid], paymentsInfo, &wg, rowsChan)
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
	header := []string{"uid", "first_name", "last_name", "email", "paypal_email", "country", "nanodegree_key", "student_uid", "week", "message_count", "activity_level", "amount"}
	suffix := fmt.Sprintf("%s_%s", startDate, endDate)
	return csv.CreateCSVFile(classMentorMessagesPaymentsFilename, suffix, header, rows)
}

func createClassMentorMessageRow(mentor mentor.Mentor, user cc.User, msgCounts map[string]map[string]map[string]int, paymentsInfo map[string]*payments.Data, wg *sync.WaitGroup, rowsChan chan<- []string) {
	defer wg.Done()

	// Create some rows.
	for ndKey, studMap := range msgCounts {
		for studUID, weekMap := range studMap {
			for week, msgCount := range weekMap {
				activityLevel, amount := displayActivityAmount(mentor.Country, msgCount, paymentsInfo[mentor.Country])
				rowsChan <- []string{
					user.ID,
					user.FirstName,
					user.LastName,
					user.Email,
					mentor.PayPalEmail,
					mentor.Country,
					ndKey,
					studUID,
					week,
					strconv.Itoa(msgCount),
					activityLevel,
					amount,
				}
			}
		}
	}
}

func calculateActivityAmount(msgCount int, classMentorMessages *payments.ClassroomMentorshipMessages) (string, float32) {
	if msgCount >= classMentorMessages.ActivityLevels["very-active"] {
		return "very-active", classMentorMessages.Amounts["very-active"]
	}
	return "active", classMentorMessages.Amounts["active"]
}

func displayActivityAmount(country string, msgCount int, paymentsData *payments.Data) (string, string) {
	if country == "" {
		return "", "COUNTRY IS NULL"
	}
	if paymentsData == nil || len(paymentsData.ClassroomMentorshipMessages.Amounts) == 0 {
		return "", fmt.Sprintf("COUNTRY '%s' IS NOT KNOWN", country)
	}
	activityLevel, amount := calculateActivityAmount(msgCount, paymentsData.ClassroomMentorshipMessages)
	return activityLevel, fmt.Sprintf("%.2f", amount)
}

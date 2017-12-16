package cm

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/udacity/mc/cc"
	"github.com/udacity/mc/mentor"
)

// CreateClassMentorAppliedRow creates a csv row for a potential class mentor.
func CreateClassMentorAppliedRow(ndKey string, currentClassMentors map[string]bool, mentor mentor.Mentor, wg *sync.WaitGroup, rowsChan chan<- []string) {
	defer wg.Done()
	// Skip if mentor is already a classroom mentor.
	if currentClassMentors[mentor.UID] {
		return
	}
	// Check if mentor applied for Classroom Mentorship.
	var appliedForClassroomMentorship bool
	for _, service := range mentor.Application.Services {
		if service == "classroom_mentorship" {
			appliedForClassroomMentorship = true
			break
		}
	}
	// Check if mentor applied for relevant Nanodegree.
	var appliedForNanodegree bool
	for _, nd := range mentor.Application.Nanodegrees {
		if nd == ndKey {
			appliedForNanodegree = true
			break
		}
	}
	if !appliedForClassroomMentorship || !appliedForNanodegree {
		return
	}

	// Fetch user. Skip if user information is not available.
	user, err := cc.FetchUserNanodegrees(mentor.UID)
	if err != nil {
		return
	}

	// Create a row.
	applicationString, err := json.Marshal(mentor.Application)
	if err != nil {
		fmt.Println("Error encoding JSON")
		return
	}

	// Convert Nanodegree representation.
	var nds []string
	for _, nd := range user.Nanodegrees {
		ndString := fmt.Sprintf("key: %s, status: %s, graduated: %t", nd.Key, nd.Enrollment.Status, nd.IsGraduated)
		nds = append(nds, ndString)
	}

	rowsChan <- []string{
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		fmt.Sprintf("%v", nds),
		mentor.PayPalEmail,
		mentor.Country,
		fmt.Sprintf("%v", mentor.Languages),
		mentor.Bio,
		mentor.EducationalBackground,
		mentor.IntroMsg,
		mentor.GitHubURL,
		mentor.LinkedInURL,
		mentor.AvatarURL,
		string(applicationString),
		mentor.CreatedAt.UTC().Format(time.RFC3339),
		mentor.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func createClassMentorByNanodegreeRow(mentor mentor.Mentor, wg *sync.WaitGroup, rowsChan chan<- []string) {
	defer wg.Done()
	// Fetch user. Skip if user information is not available.
	user, err := cc.FetchUserNanodegrees(mentor.UID)
	if err != nil {
		return
	}

	// Create a row.
	applicationString, err := json.Marshal(mentor.Application)
	if err != nil {
		fmt.Println("Error encoding JSON")
		return
	}

	// Convert Nanodegree representation.
	var nds []string
	for _, nd := range user.Nanodegrees {
		ndString := fmt.Sprintf("key: %s, status: %s, graduated: %t", nd.Key, nd.Enrollment.Status, nd.IsGraduated)
		nds = append(nds, ndString)
	}

	rowsChan <- []string{
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		fmt.Sprintf("%v", nds),
		mentor.PayPalEmail,
		mentor.Country,
		fmt.Sprintf("%v", mentor.Languages),
		mentor.Bio,
		mentor.EducationalBackground,
		mentor.IntroMsg,
		mentor.GitHubURL,
		mentor.LinkedInURL,
		mentor.AvatarURL,
		string(applicationString),
	}
}

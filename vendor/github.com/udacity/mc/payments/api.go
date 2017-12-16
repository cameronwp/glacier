package payments

import (
	"fmt"
	"net/http"

	"github.com/udacity/mc/httpclient"
)

// Response represents the response object received from Mentorship-Payments API.
type Response struct {
	Data  Data   `json:"data"`
	Error string `json:"error"`
}

// Data represents a collection of data provided by the Response.
type Data struct {
	ClassroomMentorshipMessages *ClassroomMentorshipMessages `json:"classroom-mentorship-messages,omitempty"`
	ClassroomMentorshipRatings  *ClassroomMentorshipRatings  `json:"classroom-mentorship-ratings,omitempty"`
	LiveHelpSessions            *LiveHelpSessions            `json:"live-help-sessions,omitempty"`
	Nanodegrees                 *Nanodegrees                 `json:"nanodegrees,omitempty"`
}

// ClassroomMentorshipMessages represents the classroom-mentor-message info from Mentorship-Payments API.
type ClassroomMentorshipMessages struct {
	Meta           map[string]string  `json:"meta,omitempty"`
	ActivityLevels map[string]int     `json:"activity-levels,omitempty"`
	Amounts        map[string]float32 `json:"amounts,omitempty"`
}

// ClassroomMentorshipRatings represents the classroom-mentor-rating info from Mentorship-Payments API.
type ClassroomMentorshipRatings struct {
	Meta         map[string]string  `json:"meta,omitempty"`
	Amounts      map[string]float32 `json:"amounts,omitempty"`
	RatingLevels map[string]int     `json:"rating-levels,omitempty"`
}

// LiveHelpSessions represents the LiveHelp-sessions info from Mentorship-Payments API.
type LiveHelpSessions struct {
	Meta           map[string]string             `json:"meta,omitempty"`
	Amounts        map[string]map[string]float32 `json:"amounts,omitempty"`
	ActivityLevels map[string]int                `json:"activity-levels,omitempty"`
}

// Nanodegrees represents the Nanodegree info from Mentorship-Payments API.
type Nanodegrees struct {
	Meta         map[string]string `json:"meta,omitempty"`
	Beginner     []string          `json:"beginner,omitempty"`
	Intermediate []string          `json:"intermediate,omitempty"`
	Advanced     []string          `json:"advanced,omitempty"`
}

// Fetch fetches payment info from Mentorship-Payments API using a country code.
func Fetch(isStaging bool, country string) (Data, error) {
	client, err := httpclient.New(url(isStaging))
	if err != nil {
		return Data{}, err
	}

	endpoint := fmt.Sprintf("admin/amounts?country=%s", country)
	res := Response{}
	err = client.Call(http.MethodGet, endpoint, nil, &res)
	if err != nil {
		return Data{}, err
	}

	return res.Data, nil
}

// GetNDLevel uses Data to determine the level of the supplied ndKey.
func (d *Data) GetNDLevel(ndKey string) string {
	if d != nil && d.Nanodegrees != nil && d.Nanodegrees.Beginner != nil {
		for _, currentLevelNdKey := range d.Nanodegrees.Beginner {
			if ndKey == currentLevelNdKey {
				return "beginner"
			}
		}
	}
	if d != nil && d.Nanodegrees != nil && d.Nanodegrees.Intermediate != nil {
		for _, currentLevelNdKey := range d.Nanodegrees.Intermediate {
			if ndKey == currentLevelNdKey {
				return "intermediate"
			}
		}
	}
	if d != nil && d.Nanodegrees != nil && d.Nanodegrees.Advanced != nil {
		for _, currentLevelNdKey := range d.Nanodegrees.Advanced {
			if ndKey == currentLevelNdKey {
				return "advanced"
			}
		}
	}
	return "beginner"
}

func url(isStaging bool) string {
	if isStaging {
		return stagingURL
	}
	return productionURL
}

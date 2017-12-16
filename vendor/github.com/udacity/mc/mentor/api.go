package mentor

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/udacity/mc/httpclient"
)

// Data represents a collection of data provided by the Response.
type Data struct {
	Mentor       Mentor   `json:"mentor,omitempty"`
	UpdateMentor Mentor   `json:"updateMentor,omitempty"`
	Mentors      []Mentor `json:"mentors,omitempty"`
}

// Mentor represents a mentor in Mentor-API.
type Mentor struct {
	UID                   string      `json:"uid,omitempty"`
	PayPalEmail           string      `json:"paypal_email,omitempty"`
	Country               string      `json:"country,omitempty"`
	Languages             []string    `json:"languages,omitempty"`
	Bio                   string      `json:"bio,omitempty"`
	EducationalBackground string      `json:"educational_background,omitempty"`
	IntroMsg              string      `json:"intro_msg,omitempty"`
	GitHubURL             string      `json:"github_url,omitempty"`
	LinkedInURL           string      `json:"linkedin_url,omitempty"`
	AvatarURL             string      `json:"avatar_url,omitempty"`
	Application           Application `json:"application,omitempty"`
	CreatedAt             time.Time   `json:"created_at,omitempty"`
	UpdatedAt             time.Time   `json:"updated_at,omitempty"`
}

// Application represents a raw application as sent in from mentor-dashboard.
type Application struct {
	Services    []string `json:"services,omitempty"`
	Languages   []string `json:"languages,omitempty"`
	Nanodegrees []string `json:"nanodegrees,omitempty"`
}

// FetchMentors fetches all mentors from Mentor API.
func FetchMentors(isStaging bool) ([]Mentor, error) {
	data := Data{}
	err := httpclient.PostGraphQL(httpclient.Backend{}, url(isStaging), mentorsQuery, &data)
	if err != nil {
		return nil, err
	}

	return data.Mentors, nil
}

// FetchMentor fetches a mentor from Mentor API by UID.
func FetchMentor(isStaging bool, uid string) (Mentor, error) {
	data := Data{}
	query := fmt.Sprintf(mentorQuery, uid)
	err := httpclient.PostGraphQL(httpclient.Backend{}, url(isStaging), query, &data)
	if err != nil {
		return Mentor{}, err
	}

	return data.Mentor, nil
}

// UpdateMentor will update the given fields on a mentor in mentor-api. Note:
// fields must include uid!
func UpdateMentor(isStaging bool, fields map[string]string) (Mentor, error) {
	data := Data{}
	query := buildUpdateMentorQuery(fields)
	err := httpclient.PostGraphQL(httpclient.Backend{}, url(isStaging), query, &data)
	if err != nil {
		return Mentor{}, err
	}

	return data.UpdateMentor, nil
}

func buildUpdateMentorQuery(fields map[string]string) string {
	format := func(key string, value string) string {
		return fmt.Sprintf(`%s: "%s"`, key, value)
	}

	var inputs []string

	// uid should be first
	inputs = append(inputs, format("uid", fields["uid"]))
	delete(fields, "uid")

	// anal retentively sort the rest of the keys in alphabetical order
	var sortedKeys []string
	for k := range fields {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, v := range sortedKeys {
		input := format(v, fields[v])
		inputs = append(inputs, input)
	}

	return fmt.Sprintf(mentorUpdateMutation, strings.Join(inputs, ", "))
}

func url(isStaging bool) string {
	if isStaging {
		return stagingURL
	}
	return productionURL
}

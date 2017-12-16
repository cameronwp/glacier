package guru

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/udacity/mc/cc"
	"github.com/udacity/mc/httpclient"
	"github.com/udacity/mc/mentor"
)

// Response represents the response object received from Guru.
type Response struct {
	Data     Data     `json:"data,omitempty"`
	RespGuru RespGuru `json:"guru,omitempty"`
	Errors   []Error  `json:"errors,omitempty"`
}

// Data represents a collection of data provided by the Response.
type Data struct {
	Gurus []Guru `json:"gurus,omitempty"`
}

// Error represents a collection of errors provided by the Response.
type Error struct {
	Message string `json:"message,omitempty"`
}

// Guru represents a guru in Guru.
type Guru struct {
	UID            string   `json:"uid,omitempty"`
	Name           string   `json:"name,omitempty"`
	AvatarURL      string   `json:"avatar_url,omitempty"`
	Bio            string   `json:"bio,omitempty"`
	IntroMsg       string   `json:"intro_msg,omitempty"`
	CourseIDs      []string `json:"course_ids,omitempty"`
	MaxNumStudents int      `json:"max_num_students,omitempty"`
}

// RespGuru represents a guru response in Guru (hack because course_ids is
// returned from Guru as a comma-separated string instead of an array).
type RespGuru struct {
	UID            string `json:"uid,omitempty"`
	Name           string `json:"name,omitempty"`
	AvatarURL      string `json:"avatar_url,omitempty"`
	Bio            string `json:"bio,omitempty"`
	IntroMsg       string `json:"intro_msg,omitempty"`
	CourseIDs      string `json:"course_ids,omitempty"`
	MaxNumStudents int    `json:"max_num_students,omitempty"`
	NumPairings    int    `json:"num_pairings,omitempty"`
	Status         string `json:"status,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
}

// FetchGuru sends a HTTP request to fetch a Guru from Guru.
func FetchGuru(isStaging bool, uid string) (RespGuru, error) {
	endpoint := fmt.Sprintf("api/admin/get_guru?uid=%s", uid)

	client, err := httpclient.New(url(isStaging))
	if err != nil {
		return RespGuru{}, err
	}

	res := Response{}
	err = client.Call(http.MethodPost, endpoint, nil, &res)
	if err != nil {
		return RespGuru{}, err
	}

	return res.RespGuru, nil
}

// CreateGuru creates a guru (i.e., classroom mentor) in Guru.
func CreateGuru(isStaging bool, uid string, ndkeys []string) (RespGuru, error) {
	// Fetch mentor.
	m, err := mentor.FetchMentor(isStaging, uid)
	if err != nil {
		return RespGuru{}, err
	}

	// Fetch user classroom-content.
	user, err := cc.FetchUser(uid)
	if err != nil {
		return RespGuru{}, err
	}

	// Create a new guru.
	newGuru := &Guru{
		UID:            uid,
		Name:           fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Bio:            m.Bio,
		IntroMsg:       m.IntroMsg,
		AvatarURL:      m.AvatarURL,
		CourseIDs:      ndkeys,
		MaxNumStudents: DefaultMaxNumStudents,
	}

	payload, err := json.Marshal(newGuru)
	if err != nil {
		return RespGuru{}, err
	}

	client, err := httpclient.New(url(isStaging))
	if err != nil {
		return RespGuru{}, err
	}

	res := Response{}
	err = client.Call(http.MethodPost, "/api/admin/create_guru", payload, &res)
	if err != nil {
		return RespGuru{}, err
	}

	// Fetch the newly created guru (since the /create_guru endpoint
	// used above only returns an empty JSON response).
	return FetchGuru(isStaging, uid)
}

func url(isStaging bool) string {
	if isStaging {
		return stagingURL
	}
	return productionURL
}

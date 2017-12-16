package cm

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/udacity/mc/csv"
	"github.com/udacity/mc/httpclient"
	"github.com/udacity/mc/mentor"
)

// Response represents the response object received from Classroom-Mentor API.
type Response struct {
	Data   Data    `json:"data"`
	Errors []Error `json:"errors"`
}

// Data represents a collection of data provided by the Response.
type Data struct {
	ClassMentor  ClassMentor   `json:"classmentor"`
	ClassMentors []ClassMentor `json:"classmentors"`
}

// Error represents a collection of errors provided by the Response.
type Error struct {
	Message string `json:"message"`
}

// ClassMentor represents a classroom mentor in Classroom-Mentor API.
type ClassMentor struct {
	UID            string    `json:"uid,omitempty"`
	MaxNumStudents int       `json:"max_num_students,omitempty"`
	Status         string    `json:"status,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}

// FetchClassMentor fetches a classroom mentor from Classroom-Mentor API using a UID.
func FetchClassMentor(isStaging bool, uid string) (ClassMentor, error) {
	client, err := httpclient.New(url(isStaging))
	if err != nil {
		return ClassMentor{}, err
	}

	res := Response{}
	err = client.Call(http.MethodGet, "/classmentors/"+uid, nil, &res)
	if err != nil {
		return ClassMentor{}, err
	}

	return res.Data.ClassMentor, nil
}

// FetchClassMentors fetches all classroom mentors from Classroom-Mentor API.
func FetchClassMentors(isStaging bool) ([]ClassMentor, error) {
	client, err := httpclient.New(url(isStaging))
	if err != nil {
		return []ClassMentor{}, err
	}

	res := Response{}
	err = client.Call(http.MethodGet, "/classmentors", nil, &res)
	if err != nil {
		return []ClassMentor{}, err
	}

	return res.Data.ClassMentors, nil
}

// FetchClassMentorsByNanodegree fetches classroom mentors filtered by
// Nanodegree from Classroom-Mentor API.
func FetchClassMentorsByNanodegree(isStaging bool, ndArg string) ([]ClassMentor, error) {
	client, err := httpclient.New(url(isStaging))
	if err != nil {
		return []ClassMentor{}, err
	}

	res := Response{}
	endpoint := fmt.Sprintf("classmentors?key=%s", ndArg)
	err = client.Call(http.MethodGet, endpoint, nil, &res)
	if err != nil {
		return []ClassMentor{}, err
	}

	// In case there is only one classroom mentor.
	if res.Data.ClassMentor != (ClassMentor{}) {
		return []ClassMentor{res.Data.ClassMentor}, nil
	}

	return res.Data.ClassMentors, nil
}

// FetchByND fetch current classroom mentors by Nanodegree.
func FetchByND(isStaging bool, ndkey string) (string, error) {
	classMentors, err := FetchClassMentorsByNanodegree(isStaging, ndkey)
	if err != nil {
		return "", err
	}

	// Fetch mentors.
	mentors, err := mentor.FetchMentors(isStaging)
	if err != nil {
		return "", err
	}

	mentorsMap := make(map[string]mentor.Mentor)
	for _, mentor := range mentors {
		mentorsMap[mentor.UID] = mentor
	}

	// Create rows.
	rowsChan := make(chan []string)
	var wg sync.WaitGroup
	for _, classMentor := range classMentors {
		wg.Add(1)
		go createClassMentorByNanodegreeRow(mentorsMap[classMentor.UID], &wg, rowsChan)
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
	header := []string{"uid", "first_name", "last_name", "email", "nanodegrees", "paypal_email", "country", "languages", "bio", "educational_background", "intro_msg", "github_url", "linkedin_url", "avatar_url", "application"}
	return csv.CreateCSVFile(ClassMentorsByNanodegreeFilename, ndkey, header, rows)
}

// FetchByUID fetch a mentor by UID.
func FetchByUID(isStaging bool, uid string) (ClassMentor, error) {
	cm, err := FetchClassMentor(isStaging, uid)
	if err != nil {
		return ClassMentor{}, err
	}

	return cm, nil
}

func url(isStaging bool) string {
	if isStaging {
		return stagingURL
	}
	return productionURL
}

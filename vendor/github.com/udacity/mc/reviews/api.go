package reviews

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/udacity/mc/httpclient"
)

// Certification represents a certification in reviews API.
type Certification struct {
	Status    string  `json:"status,omitempty"`
	ProjectID int     `json:"project_id,omitempty"`
	Project   Project `json:"project,omitempty"`
}

// Opportunity represents an opportunity in reviews API
type Opportunity struct {
	ID                 int       `json:"id,omitempty"`
	SubmissionRequired bool      `json:"submission_required,omitempty" csv:"submission_required,omitempty"`
	UdacityKey         string    `json:"udacity_key,omitempty" csv:"udacity_key,omitempty"`
	ProjectID          int       `json:"project_id,omitempty" csv:"project_id,omitempty"`
	Language           string    `json:"language,omitempty" csv:"language,omitempty"`
	ExpiresAt          time.Time `json:"expires_at,omitempty" csv:"expires_at,omitempty"`
	Accepted           bool      `json:"accepted,omitempty"`
	Days               int       `json:"days,omitempty" csv:"days,omitempty"`
}

// Project represents a project in reviews API.
type Project struct {
	ID             int    `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	UdacityKey     string `json:"udacity_key,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	NanodegreeKey  string `json:"nanodegree_key,omitempty"`
	AuditProjectID int    `json:"audit_project_id,omitempty"`
	AuditRubricID  int    `json:"audit_rubric_id,omitempty"`
	Error          string `json:"error,omitempty"`
}

// Submission represents a submission in reviews API.
type Submission struct {
	ID      int     `json:"id,omitempty"`
	Result  string  `json:"result,omitempty"`
	Project Project `json:"project,omitempty"`
}

// CertificationStatuses are statuses that mean a reviewer is in the certification
// flow.
var CertificationStatuses = []string{"training", "certified", "blocked"}

// FetchCerts gets certifications for a user.
func FetchCerts(isStaging bool, uid string) ([]Certification, error) {
	baseURL := url(isStaging)
	r := httpclient.Backend{}
	return fetchCerts(r, baseURL, uid)
}

func fetchCerts(r httpclient.Retriever, baseURL string, uid string) ([]Certification, error) {
	var certs []Certification
	err := httpclient.Get(r, baseURL, certsURL(uid), nil, &certs)
	if err != nil {
		return nil, err
	}
	return certs, nil
}

// FetchOpportunities gets opportunities for a user against a project.
func FetchOpportunities(isStaging bool, uid string, projectID string) ([]Opportunity, error) {
	baseURL := url(isStaging)
	r := httpclient.Backend{}
	return fetchOpportunities(r, baseURL, uid, projectID)
}

func fetchOpportunities(r httpclient.Retriever, baseURL string, uid string, projectID string) ([]Opportunity, error) {
	params := make(map[string]string)
	params["project_id"] = projectID
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	var opportunities []Opportunity
	err = httpclient.Get(r, baseURL, oppsURL(uid), payload, &opportunities)
	if err != nil {
		return nil, err
	}
	return opportunities, nil
}

// FetchProject gets project info by ID.
func FetchProject(isStaging bool, projectID string) (Project, error) {
	baseURL := url(isStaging)
	r := httpclient.Backend{}
	return fetchProject(r, baseURL, projectID)
}

func fetchProject(r httpclient.Retriever, baseURL string, projectID string) (Project, error) {
	var project Project
	endpoint := projectURL(projectID)

	err := httpclient.Get(r, baseURL, endpoint, nil, &project)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}

// FetchSubmissions gets submissions for a user against a project.
func FetchSubmissions(isStaging bool, uid string, projectID string) ([]Submission, error) {
	baseURL := url(isStaging)
	r := httpclient.Backend{}
	return fetchSubmissions(r, baseURL, uid, projectID)
}

func fetchSubmissions(r httpclient.Retriever, baseURL string, uid string, projectID string) ([]Submission, error) {
	params := make(map[string]string)
	params["project_id"] = projectID
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	var submissions []Submission
	err = httpclient.Get(r, baseURL, submissionsURL(uid), payload, &submissions)
	if err != nil {
		return nil, err
	}
	return submissions, nil
}

// PostOpportunity creates new opportunities.
func PostOpportunity(isStaging bool, opportunity Opportunity) error {
	baseURL := url(isStaging)
	r := httpclient.Backend{}
	return postOpportunity(r, baseURL, opportunity)
}

func postOpportunity(r httpclient.Retriever, baseURL string, opportunity Opportunity) error {
	params := struct {
		Opportunity Opportunity `json:"opportunity"`
	}{
		opportunity,
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return err
	}

	out := Opportunity{}
	return httpclient.Post(r, baseURL, "opportunities", payload, &out)
}

// certsURL creates a path for an API call to get all certs for a user.
func certsURL(uid string) string {
	return fmt.Sprintf("users/%s/certifications", uid)
}

// oppsURL creates a path for an API call to get all opportunities for a user.
func oppsURL(uid string) string {
	return fmt.Sprintf("users/%s/opportunities", uid)
}

// projectURL creates a path for an API call for a project ID.
func projectURL(projectID string) string {
	return fmt.Sprintf("projects/%s", projectID)
}

// submissionsURL creates a path for an API call for submissions from a user.
func submissionsURL(uid string) string {
	return fmt.Sprintf("users/by_udacity_key/%s/submissions", uid)
}

func url(isStaging bool) string {
	if isStaging {
		return stagingURL
	}
	return productionURL
}

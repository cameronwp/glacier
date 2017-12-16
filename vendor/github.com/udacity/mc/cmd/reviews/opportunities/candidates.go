package opportunities

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/udacity/mc/cc"
	"github.com/udacity/mc/csv"
	"github.com/udacity/mc/flagging"
	"github.com/udacity/mc/gae"
	"github.com/udacity/mc/lang"
	"github.com/udacity/mc/mentor"
	"github.com/udacity/mc/reviews"

	"gopkg.in/cheggaaa/pb.v1"
)

var languages = []string{}

var candidatesCmd = &cobra.Command{
	Use:   "candidates",
	Short: "Get potential candidates for a ND",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		project := reviews.Project{}
		mentors := []mentor.Mentor{}
		pc := make(chan reviews.Project)
		mc := make(chan []mentor.Mentor)
		ec := make(chan error)

		isStaging, err = cmd.Flags().GetBool("staging")
		if err != nil {
			return err
		}

		fmt.Println("Fetching project info and mentors...")
		go fetchProject(pc, ec)
		go fetchMentors(mc, ec)
		for i := 0; i < 2; i++ {
			select {
			case project = <-pc:
				break
			case mentors = <-mc:
				break
			case err = <-ec:
				if err != nil {
					return err
				}
			}
		}

		fmt.Printf("Filtering for candidates who are interested in project %d with language %s\n", project.ID, language)
		candidates := filterCandidates(mentors, project.NanodegreeKey, "reviews", languages)

		if len(candidates) == 0 {
			fmt.Println("Found 0 candidates.")
			return nil
		}

		fmt.Printf("Found %d candidate(s)\n", len(candidates))
		fmt.Println("Pulling their info (this may take a while)...")
		decorate(candidates, project.NanodegreeKey)
		candidates = pruneCertedCandidates(candidates)

		fmt.Println("Building CSV...")
		filename, err := csvify(candidates)

		fmt.Printf("Done! Results have been saved to %s\n", filename)
		return nil
	},
}

func init() {
	candidatesCmd.Flags().StringVarP(&projectID, "project_id", "p", "", "Project ID you want to find potential candidates for.")
	err := candidatesCmd.MarkFlagRequired("project_id")
	if err != nil {
		log.Fatal(err)
	}

	candidatesCmd.Flags().StringSliceVarP(&languages, "languages", "l", []string{}, "Languages you want to find potential candidates for.")
	err = candidatesCmd.MarkFlagRequired("languages")
	if err != nil {
		log.Fatal(err)
	}
}

// similar to a mentor.Mentor, with some additional fields
type candidate struct {
	UID                  string             `json:"udacity_key,omitempty"`
	Application          mentor.Application `json:"application,omitempty"`
	BlockedList          []string           `json:"blocked_list,omitempty"`
	CertificationStatus  string             `json:"certification_status,omitempty"`
	CertifiedList        []string           `json:"certified_list,omitempty"`
	CompletionAmount     float64            `json:"completion_amount,omitempty"`
	Country              string             `json:"country,omitempty"`
	Email                string             `json:"email,omitempty"`
	Flagged              string             `json:"flagged,omitempty"`
	HasOpenOpportunity   bool               `json:"has_open_opportunity,omitempty"`
	IsEnrolled           bool               `json:"is_enrolled,omitempty"`
	IsGraduated          bool               `json:"is_graduated,omitempty"`
	Languages            []string           `json:"languages,omitempty"`
	PassedProject        bool               `json:"passed_project,omitempty"`
	PrevOpportunityCount int                `json:"prev_opportunity_count,omitempty"`
	ProjectID            string             `json:"project_id,omitempty"`
	SubmissionCount      int                `json:"submission_count,omitempty"`
	CreatedAt            time.Time          `json:"created_at,omitempty"`
	UpdatedAt            time.Time          `json:"updated_at,omitempty"`
}

func fetchProject(pc chan reviews.Project, ec chan error) {
	project, err := reviews.FetchProject(isStaging, projectID)
	if err != nil {
		ec <- err
		return
	}
	pc <- project
}

func fetchMentors(mc chan []mentor.Mentor, ec chan error) {
	mentors, err := mentor.FetchMentors(isStaging)
	if err != nil {
		ec <- err
		return
	}
	mc <- mentors
}

func filterCandidates(mentors []mentor.Mentor, ndkey string, service string, languages []string) []candidate {
	candidates := filterM(mentors, func(m mentor.Mentor) bool {
		application := m.Application
		if include(application.Nanodegrees, ndkey) &&
			include(application.Services, service) &&
			includeOneOf(mapS(m.Languages, lang.Normalize), languages) {
			return true
		}
		return false
	})

	return mapToCandidates(candidates)
}

func pruneCertedCandidates(candidates []candidate) []candidate {
	return filterC(candidates, func(c candidate) bool {
		return !include(reviews.CertificationStatuses, c.CertificationStatus)
	})
}

func decorate(candidates []candidate, ndkey string) {
	count := len(candidates)
	bar := pb.StartNew(count)

	for i := range candidates {
		decorateWithCertifications(&candidates[i])

		if include(reviews.CertificationStatuses, candidates[i].CertificationStatus) {
			// don't fetch anything else, we'll remove this candidate later
			continue
		}

		decorateWithBio(&candidates[i])
		decorateWithEnrollment(&candidates[i], ndkey)
		decorateWithFlags(&candidates[i])
		decorateWithOpportunities(&candidates[i])
		decorateWithSubmissions(&candidates[i])
		bar.Increment()
	}

	bar.Finish()
}

func decorateWithBio(c *candidate) {
	user, err := gae.FetchUser(c.UID)
	if err != nil {
		log.Errorf("error: fetching gae user for %s: %s", c.UID, err)
		return
	}

	markCandidateBio(c, user)
}

func markCandidateBio(c *candidate, user gae.User) {
	c.Email = user.Email.Address
}

func decorateWithCertifications(c *candidate) {
	certs, err := reviews.FetchCerts(isStaging, c.UID)
	if err != nil {
		log.Errorf("error: fetching certifications for %s: %s", c.UID, err)
		return
	}

	markCandidateCerts(c, certs)
}

func markCandidateCerts(c *candidate, certs []reviews.Certification) {
	certificationStatus := ""
	certifiedList := []string{}
	blockedList := []string{}

	for _, cert := range certs {
		if fmt.Sprintf("%d", cert.ProjectID) == projectID {
			certificationStatus = cert.Status
		}

		if cert.Status == "certified" {
			certedProject := fmt.Sprintf("%s[%d]", cert.Project.Name, cert.Project.ID)
			certifiedList = append(certifiedList, certedProject)
		}

		if cert.Status == "blocked" {
			certedProject := fmt.Sprintf("%s[%d]", cert.Project.Name, cert.Project.ID)
			blockedList = append(blockedList, certedProject)
		}
	}

	c.CertificationStatus = certificationStatus
	c.CertifiedList = certifiedList
	c.BlockedList = blockedList
}

func decorateWithEnrollment(c *candidate, ndkey string) {
	user, err := cc.FetchUserNDProgress(c.UID, ndkey)
	if err != nil {
		log.Errorf("error: fetching Nanodegree progress for %s: %s", c.UID, err)
		return
	}

	markCandidateEnrollment(c, user)
}

func markCandidateEnrollment(c *candidate, user cc.User) {
	isEnrolled := false
	isGraduated := false
	completionAmount := 0.0

	if len(user.Nanodegrees) > 0 {
		isEnrolled = true
		isGraduated = user.Nanodegrees[0].IsGraduated
		completionAmount = user.Nanodegrees[0].AggregatedState.CompletionAmount
	}

	c.IsEnrolled = isEnrolled
	c.IsGraduated = isGraduated
	c.CompletionAmount = completionAmount
}

func decorateWithFlags(c *candidate) {
	flags, err := flagging.FetchFlags(c.UID)
	if err != nil {
		log.Errorf("error: fetching flags for %s: %s", c.UID, err)
		return
	}

	markCandidateFlags(c, flags)
}

func markCandidateFlags(c *candidate, flags []flagging.Flag) {
	flagged := "none"
	states := []string{"suspect", "cheated"}

	for _, state := range states {
		found := anyF(flags, func(f flagging.Flag) bool {
			return f.State == state
		})

		if found {
			flagged = state
		}
	}

	c.Flagged = flagged
}

func decorateWithOpportunities(c *candidate) {
	opportunities, err := reviews.FetchOpportunities(isStaging, c.UID, projectID)
	if err != nil {
		log.Errorf("error: fetching opportunities for %s | %s", c.UID, err)
		return
	}

	markCandidateOpportunities(c, opportunities)
}

func markCandidateOpportunities(c *candidate, opportunities []reviews.Opportunity) {
	c.PrevOpportunityCount = len(opportunities)
	c.HasOpenOpportunity = anyO(opportunities, func(o reviews.Opportunity) bool {
		return o.ExpiresAt.After(time.Now())
	})
}

func decorateWithSubmissions(c *candidate) {
	submissions, err := reviews.FetchSubmissions(isStaging, c.UID, projectID)
	if err != nil {
		log.Errorf("error: fetching submissions for %s: %s", c.UID, err)
		return
	}

	markCandidateSubmissions(c, submissions)
}

func markCandidateSubmissions(c *candidate, submissions []reviews.Submission) {
	c.SubmissionCount = len(submissions)
	c.PassedProject = anyS(submissions, func(s reviews.Submission) bool {
		return s.Result == "passed"
	})
}

func csvify(candidates []candidate) (string, error) {
	header, rows := buildCSVRows(candidates)
	return csv.CreateCSVFile(candidatesTemplate, projectID, header, rows)
}

func buildCSVRows(candidates []candidate) ([]string, [][]string) {
	header := []string{
		"project_id",
		"uid",
		"email",
		"country",
		"languages",
		"is_enrolled",
		"is_graduated",
		"flagged",
		"completion_amount",
		"certification_status",
		"certified_list",
		"blocked_list",
		"submission_count",
		"prev_opportunity_count",
		"has_open_opportunity",
		"passed_project",
		"created_at",
		"updated_at",
	}
	rows := [][]string{}
	for i := range candidates {
		c := &candidates[i]
		row := []string{
			c.ProjectID,
			c.UID,
			c.Email,
			c.Country,
			strings.Join(c.Languages, ","),
			fmt.Sprintf("%t", c.IsEnrolled),
			fmt.Sprintf("%t", c.IsGraduated),
			c.Flagged,
			fmt.Sprintf("%f", c.CompletionAmount),
			c.CertificationStatus,
			strings.Join(c.CertifiedList, ","),
			strings.Join(c.BlockedList, ","),
			fmt.Sprintf("%d", c.SubmissionCount),
			fmt.Sprintf("%d", c.PrevOpportunityCount),
			fmt.Sprintf("%t", c.HasOpenOpportunity),
			fmt.Sprintf("%t", c.PassedProject),
			c.CreatedAt.Local().Format(time.RFC822),
			c.UpdatedAt.Local().Format(time.RFC822),
		}
		rows = append(rows, row)
	}
	return header, rows
}

// the following are loosely from https://gobyexample.com/collection-functions

// anyS whether any submission meets a criteria.
func anyS(vs []reviews.Submission, f func(reviews.Submission) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

// anyF whether any flag meets a criteria.
func anyF(vs []flagging.Flag, f func(flagging.Flag) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

// anyO whether any opportunity meets a criteria.
func anyO(vs []reviews.Opportunity, f func(reviews.Opportunity) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

func filterM(vs []mentor.Mentor, f func(mentor.Mentor) bool) []mentor.Mentor {
	vsf := make([]mentor.Mentor, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func filterC(vs []candidate, f func(candidate) bool) []candidate {
	vsf := make([]candidate, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func include(vs []string, t string) bool {
	return index(vs, t) >= 0
}

func includeOneOf(comparator []string, comparables []string) bool {
	for _, c := range comparables {
		if include(comparator, c) {
			return true
		}
	}
	return false
}

func mapS(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func mapToCandidates(mentors []mentor.Mentor) []candidate {
	candidates := make([]candidate, len(mentors))
	for i, mentor := range mentors {
		candidates[i] = candidate{
			UID:         mentor.UID,
			Application: mentor.Application,
			Country:     mentor.Country,
			Languages:   mentor.Languages,
			CreatedAt:   mentor.CreatedAt,
			UpdatedAt:   mentor.UpdatedAt,
			ProjectID:   projectID,
		}
	}
	return candidates
}

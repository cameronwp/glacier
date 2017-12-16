package reviews

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/udacity/mc/httpclient"
)

const (
	errorRes       = `{"error":"Not Found"}`
	projectRes     = `{"id":1,"name":"Movie Trailer Website","nanodegree_key":"nd004","udacity_key":"3561209451","visible":true,"is_cert_project":false,"required_skills":"HTML, CSS, Python","audit_project_id":72,"hashtag":"nanodegree,fullstack","languages_to_recruit":[],"audit_rubric_id":1144,"awaiting_review_count":6,"awaiting_review_count_by_language":{"en-us":6},"awaiting_audit_count":0,"awaiting_quality_audit_count":0,"awaiting_training_audit_count":0,"created_at":"2015-02-03T22:44:47.603Z","updated_at":"2017-11-06T13:28:53.630Z","price":"8.82","audit_price":"7.0","rubrics":[{"id":807,"project_id":2,"description":"","upload_types":["zip","repo"],"file_filter_regex":"(^readme$)|(.*\\.(js|css|py|html|htm|txt|md|markdown|sql|swift|java|gradle|xml|rst|yml|yaml|rmd)$)","nomination_eligible":false,"stand_out":"","hide_criteria":false,"created_at":"2017-03-07T06:04:15.371Z","updated_at":"2017-03-07T06:04:15.371Z","hashtag":null,"reviewer_toolkit_url":null,"reviewer_toolkit_updated_at":null,"max_upload_size_mb":500,"canary_enabled":false,"canary_metadata":null,"estimated_sla":null,"project_assistant_enabled":false,"checkmate_enabled":false,"checkmate_metadata":null,"available_for_cert_project":false}]}`
	submissionsRes = `[{"id": 798137,"status": "canceled","result": "canceled","user_id": 72197,"notes": null,"repo_url": null,"created_at": "2017-10-15T10:14:31.880Z","updated_at": "2017-10-15T10:16:09.501Z","commit_sha": null,"grader_id": null,"assigned_at": null,"price": null,"completed_at": "2017-10-15T10:16:09.496Z","archive_url": "https://udacity-reviews-uploads.s3.amazonaws.com/_submissions/zipfile/798137/CarND-Capstone-master.zip","zipfile": {"url": "https://udacity-reviews-uploads.s3.amazonaws.com/_submissions/zipfile/798137/CarND-Capstone-master.zip"},"udacity_key": null,"held_at": null,"status_reason": null,"result_reason": null,"type": null,"training_id": null,"files": [],"url": null,"annotation_urls": [],"general_comment": "","hidden": false,"previous_submission_id": null,"nomination": null,"language": "en-us","rubric_id": 1140,"is_training": false,"canary_metadata": null,"checkmate_metadata": null,"escalated_at": null,"project_id": 331,"rubric": {"file_filter_regex": "(^readme$)|(.*\\.(js|css|py|html|htm|txt|md|markdown|sql|swift|java|gradle|xml|rst|yml|yaml|rmd)$)","upload_types": ["repo","zip"],"id": 1140,"project_id": 331,"description": "In this project you will be running student code on Carla. You can provide feedback in the project rubric along with links to datasets or videos that you want to provide to the student.","nomination_eligible": false,"stand_out": "","hide_criteria": false,"created_at": "2017-08-11T21:11:30.820Z","updated_at": "2017-10-12T04:43:54.694Z","hashtag": "","reviewer_toolkit_url": "","reviewer_toolkit_updated_at": "2017-09-27T20:01:05.183Z","max_upload_size_mb": 2000,"canary_enabled": false,"canary_metadata": null,"estimated_sla": null,"project_assistant_enabled": false,"available_for_cert_project": false,"checkmate_enabled": false,"checkmate_metadata": null},"project": {"languages_to_recruit": [],"id": 331,"name": "Programming a Real Self-Driving Car","udacity_key": "6c6fce88-d880-40f1-8b9b-df319095ec3c","created_at": "2017-08-14T17:05:24.254Z","updated_at": "2017-11-08T03:09:57.014Z","required_skills": "This will only be reviewed by Udacity employees running student code on the Udacity car, so no external reviewers are required.","visible": true,"awaiting_review_count": 65,"nanodegree_key": "nd013","audit_project_id": null,"hashtag": null,"awaiting_review_count_by_language": {"en-us": 65},"audit_rubric_id": 1144,"is_cert_project": false,"audit_price": "0.0","awaiting_audit_count": 0,"awaiting_quality_audit_count": 0,"awaiting_training_audit_count": 0,"default_price": "0.0","waitlist": true}},{"id": 785144,"status": "completed","result": "passed","user_id": 72197,"notes": null,"repo_url": "https://github.com/hristo-vrigazov/CarND-Semantic-Segmentation","created_at": "2017-10-08T20:03:33.292Z","updated_at": "2017-10-08T22:26:58.236Z","commit_sha": "7d8b90110bb4f9f348eea6c42d671f26d51c4edf","grader_id": 752,"assigned_at": "2017-10-08T20:14:03.325Z","price": "41.72","completed_at": "2017-10-08T22:26:58.234Z","archive_url": "https://udacity-reviews-uploads.s3.amazonaws.com/_submissions/zipfile/785144/archive.zip","zipfile": {"url": "https://udacity-reviews-uploads.s3.amazonaws.com/_submissions/zipfile/785144/archive.zip"},"udacity_key": null,"held_at": null,"status_reason": null,"result_reason": null,"type": null,"training_id": null,"files": [],"url": null,"annotation_urls": [],"general_comment": "Additional resources to learn more about semantic segmentation:\n\n[Stanford University | Detection and Segmentation](https://www.youtube.com/watch?v=nDPWywWRIRo)\n\n[R-CNN](https://arxiv.org/pdf/1311.2524.pdf) - Models listed below perform better; however, all papers describing them refer to this architecture.\n\n[Fully Convolutional Networks for Semantic Segmentation](https://arxiv.org/abs/1605.06211)\n\n[You Only Look Once](https://pjreddie.com/darknet/yolo/)\n\n[Single Shot MultiBox Detector](https://arxiv.org/abs/1512.02325)","hidden": false,"previous_submission_id": null,"nomination": "","language": "en-us","rubric_id": 989,"is_training": false,"canary_metadata": null,"checkmate_metadata": null,"escalated_at": null,"project_id": 311,"rubric": {"file_filter_regex": "(^readme$)|(.*\\.(js|css|py|html|htm|txt|md|markdown|sql|swift|java|gradle|xml|rst|yml|yaml|rmd)$)","upload_types": ["repo","zip"],"id": 989,"project_id": 311,"description": "# Semantic Segmentation\n### Introduction\nIn this project, you'll label the pixels of a road in images using a Fully Convolutional Network (FCN).\n\n### Setup\nClone the repo from https://github.com/udacity/CarND-Semantic-Segmentation\n##### Frameworks and Packages\nMake sure you have the following is installed:\n - [Python 3](https://www.python.org/)\n - [TensorFlow](https://www.tensorflow.org/)\n - [NumPy](http://www.numpy.org/)\n - [SciPy](https://www.scipy.org/)\n##### Dataset\nDownload the [Kitti Road dataset](http://www.cvlibs.net/datasets/kitti/eval_road.php) from [here](http://www.cvlibs.net/download.php?file=data_road.zip).","nomination_eligible": true,"stand_out": "","hide_criteria": false,"created_at": "2017-06-05T18:48:26.309Z","updated_at": "2017-07-31T17:45:56.334Z","hashtag": "","reviewer_toolkit_url": null,"reviewer_toolkit_updated_at": null,"max_upload_size_mb": 500,"canary_enabled": false,"canary_metadata": null,"estimated_sla": null,"project_assistant_enabled": false,"available_for_cert_project": false,"checkmate_enabled": false,"checkmate_metadata": null},"project": {"languages_to_recruit": [],"id": 311,"name": "Semantic Segmentation","udacity_key": "ef60b28f-3f7f-4962-bfbd-fdf82b9f31de","created_at": "2017-06-08T21:05:53.105Z","updated_at": "2017-11-08T01:49:09.071Z","required_skills": "Machine Learning\r\nPython\r\nTensorFlow\r\nConvolutional Neural Networks\r\nTransposed Convolutional Layer","visible": true,"awaiting_review_count": 0,"nanodegree_key": "nd013","audit_project_id": null,"hashtag": null,"awaiting_review_count_by_language": {},"audit_rubric_id": 1144,"is_cert_project": false,"audit_price": "33.0","awaiting_audit_count": 0,"awaiting_quality_audit_count": 0,"awaiting_training_audit_count": 0,"default_price": "35.46","waitlist": true}}]`
)

func TestFetchProject(t *testing.T) {
	projectID := "1"

	ts := httpclient.TestServer{
		Dresponse: projectRes,
	}
	testServerURL := ts.Open()
	defer ts.Close()

	retriever := httpclient.TestBackend{}
	project, err := fetchProject(retriever, fmt.Sprintf("http://%s", testServerURL), projectID)
	if err != nil {
		t.Fatalf("Something went wrong: %s", err)
	}

	stringTestCases := []struct {
		field    string
		expected string
	}{
		{"Name", "Movie Trailer Website"},
		{"UdacityKey", "3561209451"},
	}

	intTestCases := []struct {
		field    string
		expected int64
	}{
		{"ID", 1},
		{"AuditProjectID", 72},
	}

	r := reflect.ValueOf(project)

	for _, tc := range stringTestCases {
		testName := fmt.Sprintf("%s unmarshall", tc.field)
		t.Run(testName, func(t *testing.T) {
			found := reflect.Indirect(r).FieldByName(tc.field).String()
			if found != tc.expected {
				t.Errorf("got %s; expected %s", found, tc.expected)
			}
		})
	}

	for _, tc := range intTestCases {
		testName := fmt.Sprintf("%s unmarshall", tc.field)
		t.Run(testName, func(t *testing.T) {
			found := reflect.Indirect(r).FieldByName(tc.field).Int()
			if found != tc.expected {
				t.Errorf("got %d; expected %d", found, tc.expected)
			}
		})
	}
}

func TestFetchProjectWithError(t *testing.T) {
	projectID := "1"

	ts := httpclient.TestServer{
		Dresponse: errorRes,
	}
	testServerURL := ts.Open()
	defer ts.Close()

	retriever := httpclient.TestBackend{}
	project, err := fetchProject(retriever, fmt.Sprintf("http://%s", testServerURL), projectID)
	if err != nil {
		t.Fatalf("Something went wrong: %s", err)
	}
	defer ts.Close()

	if project.Error == "" {
		r, _ := json.Marshal(project)
		t.Fatalf("Expected an error in the response but did not find one: %s", r)
	}
}

func TestFetchSubmissions(t *testing.T) {
	projectID := "1"
	uid := "1"
	ts := httpclient.TestServer{
		Dresponse: submissionsRes,
	}
	testServerURL := ts.Open()
	defer ts.Close()

	retriever := httpclient.TestBackend{}
	submissions, err := fetchSubmissions(retriever, fmt.Sprintf("http://%s", testServerURL), uid, projectID)
	if err != nil {
		t.Fatalf("Something went wrong: %s", err)
	}

	if len(submissions) == 0 {
		t.Fatalf("Found 0 submissions, expected 2")
	}

	if submissions[0].ID != 798137 {
		t.Fatalf("Expected to find an ID of 798137, found %d", submissions[0].ID)
	}
}

func TestFetchSubmissionsWithError(t *testing.T) {
	projectID := "1"
	uid := "1"
	ts := httpclient.TestServer{
		Dresponse: errorRes,
	}
	testServerURL := ts.Open()
	defer ts.Close()

	retriever := httpclient.TestBackend{}
	submissions, err := fetchSubmissions(retriever, fmt.Sprintf("http://%s", testServerURL), uid, projectID)

	if err == nil {
		t.Fatalf("Something was supposed to go wrong, but didn't.")
	}

	if len(submissions) > 0 {
		t.Errorf("Expected to find no submissions, found %d instead", len(submissions))
	}
}

func TestPostOpportunity(t *testing.T) {
	k := "a12345678"
	opportunity := Opportunity{
		UdacityKey: k,
	}

	ts := httpclient.TestServer{
		ReqCB: func(r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, found %s", r.Method)
			}

			oppo := struct {
				Opportunity Opportunity `json:"opportunity"`
			}{
				Opportunity{},
			}

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
			}
			err = json.Unmarshal(body, &oppo)
			if err != nil {
				t.Error(err)
			}
			if oppo.Opportunity.UdacityKey != k {
				t.Errorf("expected the udacity key in the body to be %s, found %s", k, oppo.Opportunity.UdacityKey)
			}
		},
	}
	testServerURL := ts.Open()
	defer ts.Close()

	retriever := httpclient.TestBackend{}
	err := postOpportunity(retriever, fmt.Sprintf("http://%s", testServerURL), opportunity)
	if err != nil {
		t.Error(err)
	}
}

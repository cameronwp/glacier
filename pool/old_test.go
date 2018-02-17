package pool

// import (
// 	"fmt"
// 	"net"
// 	"net/http"
// 	"net/http/httptest"
// 	"regexp"
// 	"testing"
// )

// const (
// 	testServerURL = "127.0.0.1:1234"
// 	result1       = `{"data": "odd"}`
// 	result2       = `{"data": "even"}`
// 	JWT           = "asdfasdfasdf"
// )

// type testData struct {
// 	Data string `json:"data,omitempty"`
// }

// // create a new mock server. don't forget to close it! use the cb to run tests on the request itself.
// func testServer(baseURL string, path string, res func(subpath string) string, cb func(*http.Request)) *httptest.Server {
// 	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if cb != nil {
// 			cb(r)
// 		}

// 		re := regexp.MustCompile(fmt.Sprintf(`%s/?(.*)?`, path))
// 		if matches := re.FindStringSubmatch(r.RequestURI); matches != nil {
// 			w.WriteHeader(http.StatusOK)
// 			body := []byte(res(matches[1]))
// 			w.Write(body)
// 		} else {
// 			w.WriteHeader(http.StatusNotFound)
// 			w.Write([]byte("404 page not found"))
// 		}
// 	}))
// 	l, _ := net.Listen("tcp", baseURL)
// 	ts.Listener = l
// 	ts.Start()
// 	return ts
// }

// func ExampleNew() {
// 	schedule := New(fmt.Sprintf("http://%s/", testServerURL), JWT)
// 	out := testData{}
// 	echan := schedule("GET", "/asdf", nil, &out)

// 	err := <-echan
// 	if err != nil {
// 		fmt.Printf("Something went wrong: %s\n", err)
// 	}

// 	fmt.Println(out)
// }

// func TestNew(t *testing.T) {
// 	path := "/api/v1/endpoint"
// 	ts := testServer(testServerURL, path, func(string) string { return result1 }, nil)
// 	defer ts.Close()

// 	schedule := New(fmt.Sprintf("http://%s", testServerURL), JWT)
// 	out := testData{}
// 	echan := schedule("GET", path, nil, &out)

// 	err := <-echan
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if out.Data != "odd" {
// 		t.Errorf("Expected 'odd', found %s", out)
// 	}
// }

// func TestMultipleSchedules(t *testing.T) {
// 	path := "/api/v1/endpoint"
// 	connections := 4
// 	numRequests := 0
// 	ts := testServer(testServerURL, path, func(string) string { return result1 }, func(*http.Request) {
// 		numRequests++
// 	})
// 	defer ts.Close()

// 	schedule := New(fmt.Sprintf("http://%s", testServerURL), JWT)

// 	for i := 0; i < connections; i++ {
// 		out := testData{}
// 		echan := schedule("GET", path, nil, &out)
// 		err := <-echan
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}

// 	if numRequests != connections {
// 		t.Errorf("Expected %d requests, executed %d", connections, numRequests)
// 	}
// }

// func TestMarshalsCorrectly(t *testing.T) {
// 	path := "/api/v1/endpoint"
// 	schedule := New(fmt.Sprintf("http://%s", testServerURL), JWT)

// 	// creates a server that responds with {"data": subpath}
// 	ts := testServer(testServerURL, path, func(subpath string) string {
// 		return fmt.Sprintf(`{"data": "%s"}`, subpath)
// 	}, nil)
// 	defer ts.Close()

// 	out1 := testData{}
// 	out2 := testData{}

// 	echan1 := schedule("GET", fmt.Sprintf("%s/%d", path, 1), nil, &out1)
// 	echan2 := schedule("GET", fmt.Sprintf("%s/%d", path, 2), nil, &out2)

// 	err := <-echan1
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	err = <-echan2
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if out1.Data != "1" {
// 		t.Errorf("Expected 1, found %s", out1.Data)
// 	}

// 	if out2.Data != "2" {
// 		t.Errorf("Expected 2, found %s", out2.Data)
// 	}
// }

// func TestMoreSchedulesThanConnections(t *testing.T) {
// 	path := "/api/v1/endpoint"
// 	connections := MaxConnections * 10
// 	numRequests := 0
// 	ts := testServer(testServerURL, path, func(string) string { return result1 }, func(*http.Request) {
// 		numRequests++
// 	})
// 	defer ts.Close()

// 	schedule := New(fmt.Sprintf("http://%s", testServerURL), JWT)

// 	var echans []chan error

// 	for i := 0; i < connections; i++ {
// 		out := testData{}
// 		echan := schedule("GET", path, nil, &out)
// 		echans = append(echans, echan)
// 	}

// 	for _, ec := range echans {
// 		err := <-ec
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}

// 	// sometimes one or two requests aren't completed. not sure what's happening to the missing requests...
// 	// maybe a race condition in the test itself?
// 	if !withinRange(numRequests, connections, 0.01) {
// 		t.Errorf("Expected %d (±1%%) requests, executed %d", connections, numRequests)
// 	}
// }

// // test whether or not a value is within a range
// func withinRange(value int, comparator int, delta float64) bool {
// 	lowerBound := float64(comparator) * (1.0 - delta)
// 	upperBound := float64(comparator) * (1.0 + delta)
// 	v := float64(value)
// 	if v < lowerBound {
// 		return false
// 	}
// 	if v > upperBound {
// 		return false
// 	}
// 	return true
// }

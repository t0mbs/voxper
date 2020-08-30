package test

import (
  "testing"
  "net/http"
  "bytes"
  "io/ioutil"
)

type Case struct {
  Description string
  Endpoint string
  Method string
  Input []byte
  Status int
  Output string
}

// Test Creation using PUT
func TestPut(t *testing.T) {
  // Positive Test Cases
  cases := []Case{
    // PUT project
    Case {
      Description: "Create a new project using PUT",
      Endpoint: "/api/v1/project/1",
      Method: "PUT",
      Input: []byte(`{"name": "Project 1", "descr": "Descr for project 1"}`),
      Status: 200,
      Output: `{"status":"success","data":{"id":1,"name":"Project 1","descr":"Descr for project 1"}}`,
    },
    // PUT Engagement
    Case {
      Description: "Create a new engagement using PUT",
      Endpoint: "/api/v1/engagement/1",
      Method: "PUT",
      Input: []byte(`{"project_id": 1,"name": "Engagement 1", "descr": "Descr for engagement 1"}`),
      Status: 200,
      Output: `{"status":"success","data":{"id":1,"project_id":1,"name":"Engagement 1","descr":"Descr for engagement 1","status":"active"}}`,
    },
    // PUT Scan
    Case {
      Description: "Create a new scan using PUT",
      Endpoint: "/api/v1/scan/1",
      Method: "PUT",
      Input: []byte(`{"engagement_id": 1,"name": "Scan 1"}`),
      Status: 200,
      Output: `{"status":"success","data":{"id":1,"engagement_id":1,"name":"Scan 1","type":""}}`,
    },
    // PUT Fidning
    Case {
      Description: "Create a new finding using PUT",
      Endpoint: "/api/v1/finding/1",
      Method: "PUT",
      Input: []byte(`{"scan_id": 1,"name": "Finding 1"}`),
      Status: 200,
      Output: `{"status":"success","data":{"id":1,"scan_id":1,"name":"Finding 1","descr":"","impact":"","mitigation":"","severity":"","cve":"","cwe":"","ref":""}}`,
    },
    Case {
      Description: "Edit an existing project using PUT",
      Endpoint: "/api/v1/project/1",
      Method: "PUT",
      Input: []byte(`{"name": "Project 1 Edit"}`),
      Status: 200,
      Output: `{"status":"success","data":{"id":1,"name":"Project 1 Edit","descr":"Descr for project 1"}}`,
    },
  }
  for _, c := range cases {
    JsonCase(t, &c)
  }
}

// Test Creation using POST
func TestPost(t *testing.T) {
  // Positive Test Cases
  cases := []Case{
    // POST project
    Case {
      Description: "Create a new project using POST",
      Endpoint: "/api/v1/projects",
      Method: "POST",
      Input: []byte(`{"name": "Project 2", "descr": "Descr for project 2"}`),
      Status: 200,
      Output: `{"status":"success","data":{"id":2,"name":"Project 2","descr":"Descr for project 2"}}`,
    },
    // POST Engagement
    Case {
      Description: "Create a new engagement using POST",
      Endpoint: "/api/v1/engagements",
      Method: "POST",
      Input: []byte(`{"project_id": 2,"name": "Engagement 2", "descr": "Descr for engagement 2"}`),
      Status: 200,
      Output: `{"status":"success","data":{"id":2,"project_id":2,"name":"Engagement 2","descr":"Descr for engagement 2","status":"active"}}`,
    },
    // POST Scan
    Case {
      Description: "Create a new scan using POST",
      Endpoint: "/api/v1/scans",
      Method: "POST",
      Input: []byte(`{"engagement_id": 2,"name": "Scan 2"}`),
      Status: 200,
      Output: `{"status":"success","data":{"id":2,"engagement_id":2,"name":"Scan 2","type":""}}`,
    },
    // POST Fidning
    Case {
      Description: "Create a new finding using POST",
      Endpoint: "/api/v1/findings",
      Method: "POST",
      Input: []byte(`{"scan_id": 2,"name": "Finding 2"}`),
      Status: 200,
      Output: `{"status":"success","data":{"id":2,"scan_id":2,"name":"Finding 2","descr":"","impact":"","mitigation":"","severity":"","cve":"","cwe":"","ref":""}}`,
    },
  }
  for _, c := range cases {
    JsonCase(t, &c)
  }
}

// Test Deletion using DELETE and CASCADE rules
func TestDelete(t *testing.T) {
  // Positive Test Cases
  cases := []Case{
    // PUT project
    Case {
      Description: "Delete an existing project using DELETE and cascade",
      Endpoint: "/api/v1/project/2",
      Method: "DELETE",
      Input: []byte(""),
      Status: 200,
      Output: `{"status":"success","data":null}`,
    },
  }
  for _, c := range cases {
    JsonCase(t, &c)
  }
}

// Test Deletion and Fetching using GET
func TestGet(t *testing.T) {
  // Positive Test Cases
  cases := []Case{
    // PUT project
    Case {
      Description: "Fetch all scans using GET",
      Endpoint: "/api/v1/scans",
      Method: "GET",
      Input: []byte(""),
      Status: 200,
      Output: `{"status":"success","data":[{"id":1,"engagement_id":1,"name":"Scan 1","type":""}]}`,
    },
  }
  for _, c := range cases {
    JsonCase(t, &c)
  }
}

// Test Parsers / File Import
func TestBrakeman(t *testing.T) {

}

func TestSnyk(t *testing.T) {

}

// TODO: Implement uploads
// curl -F "engagement_id=1" -F "scan_type=Snyk" -F 'file=@test-snyk.json' -H "Content-Type: multipart/form-data" -X POST http://localhost:8000/api/v1/import
// curl -F "engagement_id=1" -F "scan_type=Brakeman" -F 'file=@test-brakeman.json' -H "Content-Type: multipart/form-data" -X POST http://localhost:8000/api/v1/import

// Helper Functions
func JsonCase (t *testing.T, c *Case) {
  req, err := http.NewRequest(c.Method, "http://localhost:8000" + c.Endpoint, bytes.NewBuffer(c.Input))
  if err != nil {
    t.Fatal(err)
  }

  client := &http.Client{}
  req.Header.Add("Content-Type", "application/json")
  resp, err := client.Do(req)
  if err != nil {
      t.Fatal(err)
  }

  defer resp.Body.Close()

  // Read response body
  bodyBytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
      t.Fatal(err)
  }

  bodyString := string(bodyBytes)

  if resp.StatusCode != c.Status {
    t.Errorf("Case failed: %v. Cause: Handler returned wrong status code: got %v want %v",
      c.Description, resp.StatusCode, c.Status)
  }

  // Compare body excluding trailing newline
  l := len(bodyString)-1
  if bodyString[:l] != c.Output {
    t.Errorf("Case failed: %v. Cause: Handler returned unexpected body: got %d %v want %d %v",
      c.Description, l, bodyString, len(c.Output), c.Output)
  }
}

// Helper Functions
func FormCase (t *testing.T, c *Case) {
  req, err := http.NewRequest(c.Method, "http://localhost:8000" + c.Endpoint, bytes.NewBuffer(c.Input))
  if err != nil {
    t.Fatal(err)
  }

  client := &http.Client{}
  req.Header.Add("Content-Type", "application/json")
  resp, err := client.Do(req)
  if err != nil {
      t.Fatal(err)
  }

  defer resp.Body.Close()

  // Read response body
  bodyBytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
      t.Fatal(err)
  }

  bodyString := string(bodyBytes)

  if resp.StatusCode != c.Status {
    t.Errorf("Case failed: %v. Cause: Handler returned wrong status code: got %v want %v",
      c.Description, resp.StatusCode, c.Status)
  }

  // Compare body excluding trailing newline
  l := len(bodyString)-1
  if bodyString[:l] != c.Output {
    t.Errorf("Case failed: %v. Cause: Handler returned unexpected body: got %d %v want %d %v",
      c.Description, l, bodyString, len(c.Output), c.Output)
  }
}

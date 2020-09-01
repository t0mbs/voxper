package main

import(
  "fmt"
  "net/http"
  "reflect"
  "encoding/json"
  "github.com/t0mbs/gorm"
  "log"
)

// Handler struct
type Handler struct {
  Database *Database
  Request *http.Request
  Writer http.ResponseWriter
  Response Response
}

// Response struct contains HTTP Response information
type Response struct {
  Status string `json:"status"`
}

type SuccessResponse struct {
    Response
    Data interface{} `json:"data"`
}

type FailData struct {
  Title string `json:"title"`
}

type FailResponse struct {
  Response
  Data FailData `json:"data"`
}

type ErrorResponse struct {
  Response
  Message string `json:"message"`
}

func (h *Handler) Close() {
  h.Request.Body.Close()
}

// Get a model from the registry by name
func (h *Handler) GetModel(name string) interface{} {
  err, model := h.Database.Registry.GetModel(name)
  if err != nil {
    h.ErrorResponse(500, err)
  }
  return model
}

// Get a slice from the registry by name
func (h *Handler) GetSlice(name string) interface{} {
  err, slice := h.Database.Registry.GetSlice(name)
  if err != nil {
    h.ErrorResponse(500, err)
  }
  return slice
}

// Unpack JSON body into a Model
func (h *Handler) UnpackBody(model interface{}) (error, interface{}) {
  decoder := json.NewDecoder(h.Request.Body)
  err := decoder.Decode(model)
  return err, model
}

// Get all records from a collection
func (h *Handler) GetCollection(name string) {
  slice := h.GetSlice(name)
  h.QueryHandler(h.Database.GetAll(slice))
}

// Get single record by id
func (h *Handler) GetRecord(id uint64, name string) {
  model := h.GetModel(name)
  h.QueryHandler(h.Database.First(id, model))
}

// Create or update database record
func (h *Handler) CreateOrUpdateRecord(id uint64, name string) {
  var err error
  // Check if record exists with URL ID
  model := h.GetModel(name)
  record := h.Database.First(id, model)

  // Unpack JSON into new model
  err, model = h.UnpackBody(model)
  if err != nil {
    h.ErrorResponse(500, err)
    return
  }

  // Set new record ID from the URL
  reflect.ValueOf(model).Elem().FieldByName("ID").SetUint(id)

  // Create record previously existed; create, if not; update
  if record.Error != nil {
    h.NewRecord(id, model)
  } else {
    h.UpdateRecord(id, model)
  }
}

func (h *Handler) UpdateRecord(id uint64, model interface{}) {
  // Save updates to record
  h.QueryHandler(h.Database.Save(model))
}

func (h *Handler) NewRecord(id uint64, model interface{}) {
  // Create new resource
  h.QueryHandler(h.Database.Create(model))

  // Setting a new record with ID desyncs the PSQL PK index, this is the fix.
  h.Database.ResetIndex(model)
}

// Create new record through POST
func (h *Handler) PostNewRecord(name string) {
  err, model := h.UnpackBody(h.GetModel(name))
  if err != nil {
    h.ErrorResponse(500, err)
  } else {
    h.QueryHandler(h.Database.Create(model))
  }
}

// Delete record and respond with a 200 or 204 depending if the record is found
func (h *Handler) DeleteRecord(id uint64, name string) {
  model := h.GetModel(name)
  response := h.Database.Delete(id, model)
  if response.Error != nil {
    h.FailResponse(204, response.Error.Error())
  } else {
    h.SuccessResponse(200, nil)
  }
}

// Import new scan
func (h *Handler) PostScan(engagement_id uint64, scan_type string, contents string) {
  // Step 1. Validate file contents as JSON
  var js map[string]interface{}
  json.Unmarshal([]byte(contents), &js)
  if js == nil {
    h.FailResponse(400, "The file supplied contains invalid JSON")
    return
  }

  // Step 2. Validate the Engagement exists
  model := h.GetModel("engagement")
  response := h.Database.First(engagement_id, model)
  if response.Error != nil {
    h.FailResponse(400, fmt.Sprintf("The Engagement with id %d does not exist", engagement_id))
    return
  }

  // Step 3. Validate the Scan Type exists
  p := Parser{}
  var findings []Finding

  switch scan_type {
  case "Snyk":
    // Create scan
    s := Scan{EngagementId: engagement_id, Type: scan_type}
    h.Database.Create(&s)
    findings = p.Snyk(s.ID, contents)

    // Create findings
    for _, f := range findings {
       h.Database.Create(&f)
    }
  case "Brakeman":
    // Create scan
    s := Scan{EngagementId: engagement_id, Type: scan_type}
    h.Database.Create(&s)
    findings = p.Brakeman(s.ID, contents)

    // Create findings
    for _, f := range findings {
       h.Database.Create(&f)
    }
  }
}

// Aggregate DB errors, respond with 500 if found, with 200 if not
func (h *Handler) QueryHandler(response *gorm.DB) {
  if response.Error != nil {
    h.ErrorResponse(500, response.Error)
    return
  } else {
    h.SuccessResponse(200, response.Value)
  }
}

func (h *Handler) SuccessResponse(code int, data interface{}) {
  // Log success
  log.Printf("[SUCCESS] Code %d, data omitted", code)
  // Send Response
  h.Respond(code, SuccessResponse{Response: Response{Status: "success"}, Data: data})
}

func (h *Handler) FailResponse(code int, title string) {
  // Log failure
  log.Printf("[FAIL] Code %d, title %s", code, title)
  // Send Response
  h.Respond(code, FailResponse{Response: Response{Status: "fail"}, Data: FailData{Title: title}})
}

// Simple Error Response
func (h *Handler) ErrorResponse(code int, err error) {
  // Log Error
  log.Printf("[ERROR] Code %d, message %s", code, err.Error())
  // Send Response
  h.Respond(code, ErrorResponse{Response: Response{Status: "error"}, Message: err.Error()})
}

// Write the HTTP response header and body
func (h *Handler) Respond(code int, response interface{}) {
    h.Writer.WriteHeader(code)
    h.Writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(h.Writer).Encode(response)
    h.Close()
}

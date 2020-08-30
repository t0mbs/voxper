package main

import (
  "net/http"
  "strconv"
  "github.com/gorilla/mux"
  "log"
  "io"
  "bytes"
)

type Router struct {
  Database *Database
}

// Init Handler
func (t *Router) Init() {
  r := mux.NewRouter()

  // Import Test Results
  r.HandleFunc("/api/v1/import", t.PostScan).Methods("POST")

  // Routes
  r.HandleFunc("/api/v1/{record:[a-z]+}s", t.GetCollection).Methods("GET")
  r.HandleFunc("/api/v1/{record:[a-z]+}s", t.PostCollection).Methods("POST")
  r.HandleFunc("/api/v1/{record:[a-z]+}/{id:[0-9]+}", t.GetRecord).Methods("GET")
  r.HandleFunc("/api/v1/{record:[a-z]+}/{id:[0-9]+}", t.PutRecord).Methods("PUT")
  r.HandleFunc("/api/v1/{record:[a-z]+}/{id:[0-9]+}", t.DeleteRecord).Methods("DELETE")

  http.Handle("/", r)
  log.Fatal(http.ListenAndServe(":8000", nil))
}

// Parse Record Args
func (t *Router) ParseUrl(r *http.Request) (uint64, string) {
  var id uint64
  vars := mux.Vars(r)
  name := vars["record"]

  // Sets ID to 0 if not specified
  if s, found := vars["id"]; found {
      id, _ = strconv.ParseUint(s, 10, 32)
  }

  return id, name
}

func (t *Router) NewHandler(w http.ResponseWriter, r *http.Request) Handler {
  return Handler{Database: t.Database, Request: r, Writer: w}
}

// API Callbacks
func (t *Router) PostScan(w http.ResponseWriter, r *http.Request){
  h := t.NewHandler(w, r)
  engagementId, _ := strconv.ParseUint(r.FormValue("engagement_id"), 10, 32)
  scanType := r.FormValue("scan_type")

  // Max input size of 1*2^20 Bytes, or 10MB
  r.ParseMultipartForm(1 << 20)
  var buf bytes.Buffer

  // Open file
  file, _, err := r.FormFile("file")
  if err != nil {
    h.FailResponse(400, err.Error())
    return
  }
  defer file.Close()

  // Copy the file data into buffer
  io.Copy(&buf, file)
  contents := buf.String()
  buf.Reset()

  // Import contents of JSON file
  h.PostScan(engagementId, scanType, contents)
}

func (t *Router) GetCollection(w http.ResponseWriter, r *http.Request){
  _, name := t.ParseUrl(r)
  h := t.NewHandler(w, r)
  h.GetCollection(name)
}

func (t *Router) PostCollection(w http.ResponseWriter, r *http.Request){
  _, name := t.ParseUrl(r)
  h := t.NewHandler(w, r)
  h.PostNewRecord(name)
}

func (t *Router) GetRecord(w http.ResponseWriter, r *http.Request) {
  id, name := t.ParseUrl(r)
  h := t.NewHandler(w, r)
  h.GetRecord(id, name)
}

func (t *Router) PutRecord(w http.ResponseWriter, r *http.Request) {
  id, name := t.ParseUrl(r)
  h := t.NewHandler(w, r)
  h.CreateOrUpdateRecord(id, name)
}

func (t *Router) DeleteRecord(w http.ResponseWriter, r *http.Request) {
  id, name := t.ParseUrl(r)
  h := t.NewHandler(w, r)
  h.DeleteRecord(id, name)
}

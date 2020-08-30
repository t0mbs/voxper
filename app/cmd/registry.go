package main

import (
  "time"
  "errors"
  "reflect"
  "strings"
)

// Add models to slice to populate Registry, tables will be created in order
var models = []interface{} {Project{}, Engagement{}, Scan{}, Finding{}}

type Registry struct {
  Models map[string]interface{}
}

func (r *Registry) Populate() {
  // Initiate empty map
  r.Models = make(map[string]interface{})
  for _, model := range models {
    // Get string of Type name
    name := reflect.ValueOf(model).Type().Name()
    // To Lowercase
    lc := strings.ToLower(name)
    r.Models[lc] = model
  }
}

func (r *Registry) GetModel(name string) (error, interface{}) {
  // If model exists
  if model, ok := r.Models[name]; ok {
    // Get type of model
    t := reflect.TypeOf(model)
    // New model
		m  := reflect.New(t).Interface()
    // Return model pointer
    return nil, m
  }

  return errors.New("Model not found"), nil
}

func (r *Registry) GetSlice(name string) (error, interface{}) {
  // If model exists
  if model, ok := r.Models[name]; ok {
    // Get type of struct
    t := reflect.TypeOf(model)
    // New slice of structs
		s := reflect.MakeSlice(reflect.SliceOf(t), 0, 0).Interface()
    // New assignable slice of structs
    i := reflect.New(reflect.TypeOf(s)).Interface()
    // Return slice pointer
    return nil, i
  }

  return errors.New("Model not found"), nil
}

//  Generic model definitions
type Model struct {
	ID        uint64       `gorm:"primary_key" json:"id"`
  CreatedAt time.Time    `json:"-"`
  UpdatedAt time.Time    `json:"-"`
}

// Project represents to a single product or repository
type Project struct {
  Model
  Name    string `gorm:"unique;not null" json:"name"`
  Descr   string `json:"descr"`
}

// Engagement represents a grouping of tests, run manually or by CI/CD pipeline
type Engagement struct {
  Model
  Project Project `gorm:"foreignkey:ProjectId" json:"-"`
  ProjectId uint64 `gorm:"not null" json:"project_id"`
  Name string `gorm:"not null" json:"name"`
  Descr string `json:"descr"`
  Status string `json:"status" gorm:"default:'active'"`
}

// Scan represents an automatic or manual test
type Scan struct {
  Model
  Engagement Engagement `gorm:"foreignkey:EngagementId" json:"-"`
  EngagementId uint64 `gorm:"not null" json:"engagement_id"`
  Name string `json:"name"`
  Type string `gorm:"not null" json:"type"`
}

// Finding is a finding, confirmed or unconfirmed
type Finding struct {
  Model
  Scan Scan `gorm:"foreignkey:ScanId" json:"-"`
  ScanId uint64 `gorm:"not null" json:"scan_id"`
  Name string `gorm:"not null" json:"name"`
  Descr string `json:"descr"`
  Impact string `json:"impact"`
  Mitigation string `json:"mitigation"`
  Severity string `json:"severity"`
  Cve string `json:"cve"`
  Cwe string `json:"cwe"`
  Ref string `json:"ref"`
}

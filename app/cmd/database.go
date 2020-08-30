package main

import (
  "fmt"
  "log"
  "reflect"
  "strings"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

// Struct holding Go-ORM database pointer, registry and connection information
type Database struct {
  Connector *gorm.DB
  Registry *Registry
  Host string
  Port string
  Name string
  User string
  Pass string
}

// Init database connection and create tables
func (db *Database) Init() {
  var err error
  c := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", db.Host, db.Port, db.Name, db.User, db.Pass)
  db.Connector, err = gorm.Open("postgres", c)

  if err != nil {
    log.Fatal(err)
  }

  // Create tables
  for _, model := range models {
    if !db.Connector.HasTable(model) {
      log.Printf("Creating table %T", model)
      db.Connector.CreateTable(model)

      // Create foreign key associations, not handled by GORM default
      t := reflect.TypeOf(model)
      for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        tag := field.Tag.Get("gorm")
        // If field has a defined GORM foreignkey
        if strings.Contains(tag, "foreignkey:") {
          // Get foreignkey name
          foreignKey := strings.Split(tag, ":")[1]
          // Get GORM column name for foreign key column
          col := gorm.ToColumnName(foreignKey)
          // Get reference column, assuming id primary key
          refCol := gorm.ToTableName(field.Name) + "s(id)"
          // Add foreignkey constraint with cascade on update / delete
          db.Connector.Model(model).AddForeignKey(col, refCol, "CASCADE", "CASCADE")
        }
      }
    }
  }
}

// Close database connection
func (db *Database) Close() {
  db.Connector.Close()
}

// Get all records in table
func (db *Database) GetAll(slice interface{}) *gorm.DB {
  return db.Connector.Find(slice)
}

// Select first record by ID
func (db *Database) First(id uint64, model interface{}) *gorm.DB {
  return db.Connector.First(model, id)
}

// Create record
func (db *Database) Create(model interface{}) *gorm.DB {
  return db.Connector.Create(model)
}

// Save record
func (db *Database) Save(model interface{}) *gorm.DB {
  return db.Connector.Save(model)
}

// Delete record by id
func (db *Database) Delete(id uint64, model interface{}) *gorm.DB {
  return db.Connector.Unscoped().Where("id = ?", id).Delete(model)
}

// Reset Primary Key index for table
func (db *Database) ResetIndex(model interface{}) {
  table := db.Connector.NewScope(model).TableName()
  db.Connector.Exec(fmt.Sprintf("SELECT setval('%s_id_seq', max(id)) FROM %s", table, table))
}

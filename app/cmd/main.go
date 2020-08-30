package main

import "os"

func main() {
  // Init Registry
  r := Registry{}
  r.Populate()

  // Register models with Database
  db := Database{Registry: &r, Host: os.Getenv("DB_HOST"), Port: os.Getenv("DB_PORT"), Name: os.Getenv("DB_NAME"), User: os.Getenv("DB_USER"), Pass: os.Getenv("DB_PASS")}

  // Init Database Connection
  db.Init()
  defer db.Close()

  // Init Server / Router
  t := Router{Database: &db}
  t.Init()
}

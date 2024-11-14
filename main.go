package main

import (
	_ "embed"
	"encoding/json"
	"log"

	"github.com/rus-sharafiev/go-push/db"
	"github.com/rus-sharafiev/go-push/push"
)

//go:embed config.json
var config []byte

func main() {
	if err := json.Unmarshal(config, &push.Config); err != nil {
		log.Fatalf("\n\x1b[31m Error parsing the config file: %v\x1b[0m\n", err)
	}

	if err := db.Connect(*push.Config.Db); err != nil {
		log.Fatalf("\n\x1b[31m Error connecting the database: %v\x1b[0m\n", err)
	}

	log.Fatal(push.Service.Start())
}

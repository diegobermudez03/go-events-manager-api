package main

import (
	"log"
	"os"

	"github.com/diegobermudez03/go-events-manager-api/internal/api"
	"github.com/diegobermudez03/go-events-manager-api/internal/config"
	"github.com/diegobermudez03/go-events-manager-api/internal/db"
	"github.com/diegobermudez03/go-events-manager-api/pkg/storage"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

//My strategy is that the API shoulnd't have to be responsible of starting services, like databases
//redis, etc, it should have everything injected, and it is only responsible of the dependency injection
//and running the server, which means, all the external services setup is made by the main, the API
//shouldnt have to deal with errors with those services, it receives the already opened services, as done with Storage

func main() {
	//load env variables and configuration
	godotenv.Load(".env")
	config := &config.Config{
		Port: getEnv("PORT", ":8081"),
		DbConfig: config.DbConfig{
			Addr: getEnv("POSTGRES_URL", "postgres://admin:secret@localhost:5432/events_go?sslmode=disable"),
		},
	}

	//open database
	db, err := db.NewDatabase(config.DbConfig.Addr)
	if err != nil{
		log.Fatalf("Unable to open database %s", err.Error())
	}
	defer db.Close()

	storage := storage.NewPostgreStorage(db)

	//MIGRATIONS
	m, err := migrate.New(
		"file://cmd/migrations/migrate",
		config.DbConfig.Addr,
	)
	if err != nil{
		log.Fatalf("Unable to migrate %s", err.Error())
	}
	if err = m.Up(); err != nil{
		log.Fatalf("Unable to migrate %s", err.Error())
	}
	log.Println("Migrations up succesfully")

	//create new API server
	server := api.NewAPIServer(storage)

	//run the API server
	if err := server.Run(); err != nil{
		log.Fatalf("couldn't start server %s", err.Error())
	}
}


func getEnv(param string, fallback string) string{
	if val, ok := os.LookupEnv(param); ok{
		return val
	}
	return fallback
}
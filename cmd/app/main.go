package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/diegobermudez03/go-events-manager-api/internal/config"
	"github.com/diegobermudez03/go-events-manager-api/internal/db"
	"github.com/diegobermudez03/go-events-manager-api/internal/http/api"
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
	config := config.NewConfig()

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
	if err = m.Up(); err != nil && err != migrate.ErrNoChange{
		log.Fatalf("Unable to migrate %s", err.Error())
	}
	log.Println("Migrations up succesfully")

	//for graceful shutdown, listen to os signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	//create new API server
	server := api.NewAPIServer(config.Port, storage, config)
	//run the API server in a separated gorutine to be able to listen to shutdown
	go func(){
		log.Printf("Server running on port %s", config.Port)
		if err := server.Run(); err != nil && err != http.ErrServerClosed{
			log.Fatalf("couldn't start server %s", err.Error())
		}
	}()

	//waiting for the stop signal, is just a blocking code
	<-ctx.Done()
	log.Println("Interruption signal")
	if err := server.Shutdown(); err != nil{
		log.Fatalf("Server shutdown error %s", err.Error())
	}
	log.Println("Succesfully graceful shutdown")
}


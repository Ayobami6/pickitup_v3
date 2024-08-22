package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Ayobami6/pickitup_v3/cmd/api"
	"github.com/Ayobami6/pickitup_v3/config"
	"github.com/Ayobami6/pickitup_v3/db"
	_ "github.com/joho/godotenv/autoload"
)

func init() {
	file, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
        log.Fatalf("Failed to open log file: %v", err)
    }
	multiWriter := io.MultiWriter(file, os.Stdout)
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	host := config.GetEnv("DB_HOST", "localhost")
	port := config.GetEnv("DB_PORT", "5432")
	user := config.GetEnv("DB_USER", "ayo")
	pwd := config.GetEnv("DB_PWD", "password")
	dbName := config.GetEnv("DB_NAME", "pickitup_db")
	sslmode := "disable"
    timeZone := "Africa/Lagos"
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s Timezone=%s", user, pwd, dbName, host, port, sslmode, timeZone)
	DB, err := db.ConnectDb(dsn)
	if err!= nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
	addr := config.GetEnv("ADDRESS_URL", "localhost:2300")
	apiServer := api.NewAPIServer(addr, DB)
	if err := apiServer.Run(); err!= nil {
        log.Fatal(err)
    }
    
}
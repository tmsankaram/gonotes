package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/tmsankram/gonotes/internal/config"
	"github.com/tmsankram/gonotes/internal/db"
	"github.com/tmsankram/gonotes/internal/notes"
	"github.com/tmsankram/gonotes/internal/router"
	"github.com/tmsankram/gonotes/internal/users"
)

// @title GoNotes API
// @version 1.0
// @description GoNotes — Notes API with Gin, GORM, OAuth, TOTP, GraphQL (read).
// @termsOfService http://swagger.io/terms/

// @contact.name Deva
// @contact.email you@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()  // load configuration
	db := db.Connect(cfg) // connect to the database

	// AutoMigrate models
	if err := db.AutoMigrate(&users.User{}, &notes.Note{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	r := router.New(db)

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Port),
		Handler: r,
	}

	go func() {
		log.Printf("Starting GoNotes on :%d", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server Error: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down...")
	srv.Shutdown(ctx)
}

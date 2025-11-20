package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"github.com/tmsankram/gonotes/internal/auth"
	"github.com/tmsankram/gonotes/internal/files"
	"github.com/tmsankram/gonotes/internal/middleware"
	"github.com/tmsankram/gonotes/internal/notes"
	"github.com/tmsankram/gonotes/internal/users"
	myval "github.com/tmsankram/gonotes/internal/validator"
)

func New(db *gorm.DB) *gin.Engine {
	r := gin.New()

	// Global Middleware
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// register custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("notest", myval.TitleNoTest)
	}

	// Health check
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	// Services
	notesSvc := notes.NewService(db) // notes service depends on db
	usersSvc := users.NewService(db) // users service depends on db
	filesSvc := files.NewService()   // files service using in memory storage

	// Notes Handler
	notesHandler := notes.NewHandler(notesSvc)
	notesHandler.RegisterRoutes(r)

	// Files Handler
	filesHandler := files.NewHandler(filesSvc)
	filesHandler.RegisterRoutes(r)

	authHandler := auth.NewHandler(usersSvc)
	totpHandler := auth.NewTOTPHandler(usersSvc)

	// Auth routes moved to handlers
	authHandler.RegisterPublicRoutes(r)
	authHandler.RegisterProtectedRoutes(r)
	totpHandler.RegisterRoutes(r)

	return r
}

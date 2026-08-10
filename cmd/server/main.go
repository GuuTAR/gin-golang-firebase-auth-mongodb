// @title           gin-golang-firebase-auth-mongodb API
// @version         1.0
// @description     Reusable Go project template using Gin, Firebase Authentication, and MongoDB.
//
// @BasePath        /
//
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Firebase ID token — prefix with "Bearer "
package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/GuuTAR/gin-golang-firebase-auth-mongodb/docs"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/config"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/handlers"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/middleware"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/pkg/auth"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/pkg/db"
)

func main() {
	// Load .env for local development; ignore error when file is absent.
	_ = godotenv.Load()

	cfg := config.Load()
	// ── Process-lifetime context (cancelled on SIGINT / SIGTERM) ─────────────
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// ── MongoDB ──────────────────────────────────────────────────────────────
	mongoClient, err := db.NewMongoClient(ctx, cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatalf("mongodb: %v", err)
	}
	log.Printf("mongodb connected (db=%s)", cfg.MongoDB)

	// ── Firebase ─────────────────────────────────────────────────────────────
	firebaseClient, err := auth.NewFirebaseClient(ctx, cfg.Firebase)
	if err != nil {
		log.Fatalf("firebase: %v", err)
	}
	log.Println("firebase auth client initialised")

	// ── Dependency wiring ────────────────────────────────────────────────────
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler()
	tokenHandler := handlers.NewTokenHandler(cfg.FirebaseWebAPIKey)

	// ── Router ───────────────────────────────────────────────────────────────
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Public routes
	r.GET("/health", healthHandler.Check)
	r.GET("/docs", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/docs/index.html") })
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")

		// POST /api/v1/auth/token — sign in with email+password, returns a Bearer token.
		// Useful for testing with curl / Postman. In production the client app
		// should call Firebase directly using the client SDK.
		authGroup.POST("/token", tokenHandler.Token)

		// Protected routes — require a valid Firebase ID token
		authGroup.Use(middleware.FirebaseAuth(firebaseClient.Auth))
		{
			authGroup.GET("/me", authHandler.Me)
		}
	}

	// ── HTTP server ──────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("server listening on :%s (env=%s)", cfg.Port, cfg.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// ── Graceful shutdown ────────────────────────────────────────────────────
	<-ctx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}
	if err := mongoClient.Disconnect(shutdownCtx); err != nil {
		log.Printf("mongodb disconnect: %v", err)
	}
	log.Println("server exited")
}

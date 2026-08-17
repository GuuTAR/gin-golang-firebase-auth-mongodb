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

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/GuuTAR/gin-golang-firebase-auth-mongodb/docs"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/config"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/handlers"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/middleware"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/repository"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/services"
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
	authSvc := services.NewAuthService()
	authHandler := handlers.NewAuthHandler(cfg.FirebaseWebAPIKey, authSvc)

	todoRepo := repository.NewTodoRepository(mongoClient.DB)
	todoSvc := services.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoSvc)

	// ── Router ───────────────────────────────────────────────────────────────
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Public routes
	r.GET("/health", healthHandler.Check)
	r.GET("/docs", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/docs/index.html") })
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")

		// Protected routes — require a valid Firebase ID token
		authGroup.Use(middleware.FirebaseAuth(firebaseClient.Auth))
		{
			authGroup.GET("/me", authHandler.Me)
		}
		v1.POST("/auth/signin", authHandler.SignIn)

		// Todo routes — all protected by Firebase auth
		todoGroup := v1.Group("/todos")
		todoGroup.Use(middleware.FirebaseAuth(firebaseClient.Auth))
		{
			todoGroup.POST("", todoHandler.Create)
			todoGroup.GET("", todoHandler.List)
			todoGroup.DELETE("/:id", todoHandler.Delete)
		}
		v1.PATCH("/todos-reset", middleware.FirebaseAuth(firebaseClient.Auth), todoHandler.Reset)
		v1.PATCH("/todos-toggle/:id", middleware.FirebaseAuth(firebaseClient.Auth), todoHandler.Toggle)
	}

	// ── HTTP server ──────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		// Longer than any fronting proxy/tunnel's idle timeout, so it never closes
		// a keep-alive connection the proxy still considers live and reuses.
		IdleTimeout: 120 * time.Second,
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

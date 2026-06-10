package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"

	_ "music-service/docs"
	"music-service/internal/core/config"
	"music-service/internal/core/logger"
	"music-service/internal/core/middleware"
	"music-service/internal/core/postgres"
	"music-service/internal/core/redis"
	auth_handler "music-service/internal/features/auth/handler"
	auth_postgres "music-service/internal/features/auth/repository/postgres"
	auth_service "music-service/internal/features/auth/service"
	favorites_handler "music-service/internal/features/favorites/handler"
	favorites_postgres "music-service/internal/features/favorites/repository/postgres"
	favorites_service "music-service/internal/features/favorites/service"
	history_handler "music-service/internal/features/history/handler"
	history_postgres "music-service/internal/features/history/repository/postgres"
	playlists_handler "music-service/internal/features/playlists/handler"
	playlists_postgres "music-service/internal/features/playlists/repository/postgres"
	playlists_service "music-service/internal/features/playlists/service"
	subscriptions_handler "music-service/internal/features/subscriptions/handler"
	tracks_handler "music-service/internal/features/tracks/handler"
	tracks_postgres "music-service/internal/features/tracks/repository/postgres"
	tracks_service "music-service/internal/features/tracks/service"
	users_handler "music-service/internal/features/users/handler"
	users_postgres "music-service/internal/features/users/repository/postgres"
	users_service "music-service/internal/features/users/service"
)

// @title           Music Subscription Service API
// @version         1.0
// @description     This is a music subscription service server.
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     Введите токен в формате: Bearer <ваш_токен>
func main() {
	// 1. Load configuration
	cfg := config.Load()

	// 2. Initialize Logger
	log, err := logger.New(cfg.Env)
	if err != nil {
		fmt.Printf("failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("starting application", zap.String("env", cfg.Env))

	// 3. Connect to Postgres
	var db *sqlx.DB
	if cfg.DatabaseURL != "" {
		db, err = postgres.NewConnectionFromURL(cfg.DatabaseURL)
	} else {
		db, err = postgres.NewConnection(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
	}
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()
	log.Info("connected to postgres database")

	// 4. Connect to Redis
	rdb, err := redis.NewClient(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword)
	if err != nil {
		log.Fatal("failed to connect to redis", zap.Error(err))
	}
	defer rdb.Close()
	log.Info("connected to redis")

	// 4.1 Initialize Cache Client
	cache := redis.NewCache(rdb)

	// 5. Dependency Injection (DI)
	
	// Users & Auth
	userRepo := users_postgres.NewUserRepository(db)
	userService := users_service.NewUserService(userRepo)
	userHandler := users_handler.NewUserHandler(userService)

	authRefreshRepo := auth_postgres.NewRefreshTokenRepository(db)
	authService := auth_service.NewAuthService(
		userRepo,
		authRefreshRepo,
		cfg.JWTAccessSecret,
		cfg.JWTRefreshSecret,
		cfg.JWTAccessTTL,
		cfg.JWTRefreshTTL,
	)
	authHandler := auth_handler.NewAuthHandler(authService)

	// Tracks
	trackRepo := tracks_postgres.NewTrackRepository(db)
	trackService := tracks_service.NewTrackService(trackRepo, cache, rdb, cfg.FreeDailyPlayLimit)
	trackHandler := tracks_handler.NewTrackHandler(trackService)

	// Playlists
	playlistRepo := playlists_postgres.NewPlaylistRepository(db)
	playlistService := playlists_service.NewPlaylistService(playlistRepo, cfg.FreePlaylistLimit)
	playlistHandler := playlists_handler.NewPlaylistHandler(playlistService)

	// Favorites
	favoritesRepo := favorites_postgres.NewFavoritesRepository(db)
	favoritesService := favorites_service.NewFavoritesService(favoritesRepo, 20) // 20 favorites limit for FREE
	favoritesHandler := favorites_handler.NewFavoritesHandler(favoritesService)

	// Listening History
	historyRepo := history_postgres.NewHistoryRepository(db)
	historyHandler := history_handler.NewHistoryHandler(historyRepo)

	// Subscriptions
	subscriptionHandler := subscriptions_handler.NewSubscriptionHandler(userService)

	// 6. Router Setup
	r := chi.NewRouter()

	// Global Middlewares
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(log))
	r.Use(middleware.Trace())
	r.Use(middleware.Recoverer())

	// Public Health Route
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Swagger Route
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// API Routing
	r.Route("/api/v1", func(r chi.Router) {
		// Public Auth Routes
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.Refresh)
		r.Post("/auth/logout", authHandler.Logout)

		// Protected Routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg.JWTAccessSecret))

			// Users Profile
			r.Get("/users/me", userHandler.GetMe)
			r.Patch("/users/me", userHandler.UpdateProfile)

			// Catalog View & Play
			r.Get("/tracks", trackHandler.ListTracks)
			r.Get("/tracks/{id}", trackHandler.GetTrack)
			r.Get("/tracks/search", trackHandler.SearchTracks)
			r.Post("/tracks/{id}/play", trackHandler.PlayTrack)

			// Playlists
			r.Get("/playlists", playlistHandler.ListPlaylists)
			r.Post("/playlists", playlistHandler.CreatePlaylist)
			r.Get("/playlists/{id}", playlistHandler.GetPlaylist)
			r.Put("/playlists/{id}", playlistHandler.UpdatePlaylist)
			r.Delete("/playlists/{id}", playlistHandler.DeletePlaylist)
			r.Post("/playlists/{playlist_id}/tracks/{track_id}", playlistHandler.AddTrack)
			r.Delete("/playlists/{playlist_id}/tracks/{track_id}", playlistHandler.RemoveTrack)

			// Favorites
			r.Get("/favorites/tracks", favoritesHandler.ListFavorites)
			r.Post("/favorites/tracks/{track_id}", favoritesHandler.AddFavorite)
			r.Delete("/favorites/tracks/{track_id}", favoritesHandler.RemoveFavorite)

			// History
			r.Get("/listening-history", historyHandler.GetHistory)

			// Admin endpoints group
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole("ADMIN"))

				// Tracks Management
				r.Post("/admin/tracks", trackHandler.CreateTrack)
				r.Put("/admin/tracks/{id}", trackHandler.UpdateTrack)
				r.Delete("/admin/tracks/{id}", trackHandler.DeleteTrack)

				// Subscriptions Modification
				r.Patch("/admin/users/{id}/subscription", subscriptionHandler.UpdateSubscription)
			})
		})
	})

	// 7. Start HTTP Server with Graceful Shutdown
	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("server is running", zap.String("port", cfg.AppPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed to listen", zap.Error(err))
		}
	}()

	// Signal channels
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("shutting down HTTP server...")

	// Context for graceful shutdown (timeout: 10 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", zap.Error(err))
	} else {
		log.Info("server stopped gracefully")
	}
}


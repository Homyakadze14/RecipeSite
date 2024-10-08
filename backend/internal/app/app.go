// Package app configures and runs application.
package app

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Homyakadze14/RecipeSite/config"
	v1 "github.com/Homyakadze14/RecipeSite/internal/controller/http/v1"
	"github.com/Homyakadze14/RecipeSite/internal/filestorage"
	repo "github.com/Homyakadze14/RecipeSite/internal/repository/postgres"
	rabbitmqrepo "github.com/Homyakadze14/RecipeSite/internal/repository/rabbitmq"
	redisrepo "github.com/Homyakadze14/RecipeSite/internal/repository/redis"
	"github.com/Homyakadze14/RecipeSite/internal/usecases"
	"github.com/Homyakadze14/RecipeSite/pkg/httpserver"
	"github.com/Homyakadze14/RecipeSite/pkg/postgres"
	"github.com/Homyakadze14/RecipeSite/pkg/rabbitmq"
	"github.com/Homyakadze14/RecipeSite/pkg/redis"
	"github.com/gin-gonic/gin"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		slog.Error(fmt.Errorf("app - Run - postgres.New: %w", err).Error())
		os.Exit(1)
	}
	defer pg.Close()

	s3, err := filestorage.NewS3Storage(cfg)
	if err != nil {
		slog.Error(fmt.Errorf("app - Run - filestorage.NewS3Storage: %w", err).Error())
		os.Exit(1)
	}

	// RMQ
	rmq, err := rabbitmq.New(cfg.RMQ.URL)
	if err != nil {
		slog.Error(fmt.Errorf("app - Run - rabbitmq.New: %w", err).Error())
		os.Exit(1)
	}
	defer rmq.Close()

	// Redis
	redis, err := redis.New(cfg.Redis)
	if err != nil {
		slog.Error(fmt.Errorf("app - Run - redis.New: %w", err).Error())
		os.Exit(1)
	}
	defer redis.Close()

	redisRepo := redisrepo.NewRedisRepository(redis)

	// Rabbit repository
	rmqRepo, err := rabbitmqrepo.NewSubscribeRabbitMQRepository(rmq)
	if err != nil {
		slog.Error(fmt.Errorf("app - Run - rabbitmqrepo.NewSubscribeRabbitMQRepository: %w", err).Error())
		os.Exit(1)
	}
	defer rmqRepo.CloseChan()

	// Use cases
	sessionUseCase := usecases.NewSessionUseCase(repo.NewSessionRepository(pg))
	likeUseCase := usecases.NewLikeUsecase(repo.NewLikeRepository(pg))
	jwtUseCase := usecases.NewJWTUseCase([]byte(cfg.JWT.SECRET_KEY))
	userUseCase := usecases.NewUserUsecase(repo.NewUserRepository(pg), sessionUseCase, cfg.DEFAULT_ICON_URL, s3, jwtUseCase, redisRepo, likeUseCase)
	commentUseCase := usecases.NewCommentUseCase(repo.NewCommentRepository(pg), userUseCase)
	subscribeUseCase := usecases.NewSubscribeUsecase(repo.NewSubscribeRepository(pg), rmqRepo, userUseCase)
	recipeUseCase := usecases.NewRecipeUsecase(repo.NewRecipeRepository(pg), userUseCase, likeUseCase,
		s3, commentUseCase, subscribeUseCase, redisRepo)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, sessionUseCase, userUseCase, likeUseCase, recipeUseCase, commentUseCase, subscribeUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	slog.Info("app - Run - start working")

	select {
	case s := <-interrupt:
		slog.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		slog.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err).Error())

		// Shutdown
		err = httpServer.Shutdown()
		if err != nil {
			slog.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err).Error())
		}
	}
}

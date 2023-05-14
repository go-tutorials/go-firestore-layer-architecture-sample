package app

import (
	"context"
	"firebase.google.com/go"
	v "github.com/core-go/core/v10"
	"github.com/core-go/health"
	"github.com/core-go/health/firestore"
	"github.com/core-go/log"
	"google.golang.org/api/option"

	"go-service/internal/handler"
	"go-service/internal/repository"
)

type ApplicationContext struct {
	Health *health.Handler
	User   handler.UserPort
}

func NewApp(ctx context.Context, cfg Config) (*ApplicationContext, error) {
	opts := option.WithCredentialsJSON([]byte(cfg.Credentials))
	app, err := firebase.NewApp(ctx, nil, opts)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	logError := log.LogError
	validator := v.NewValidator()

	userService := repository.NewUserRepository(client)
	userHandler := handler.NewUserHandler(userService, validator.Validate, logError)

	firestoreChecker := firestore.NewHealthChecker(ctx, []byte(cfg.Credentials))
	healthHandler := health.NewHandler(firestoreChecker)

	return &ApplicationContext{
		Health: healthHandler,
		User:   userHandler,
	}, nil
}

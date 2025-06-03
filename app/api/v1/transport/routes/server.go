package routes

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"

	v1 "github.com/me0den/example-service/app/api/v1"
)

// ServerFXModule represents a FX module for http server.
var ServerFXModule = fx.Options(
	fx.Invoke(
		startHTTPServer,
	),
)

// Validator custom validator for echo.
type Validator struct {
	Validator *validator.Validate
}

// Validate implement custom validate.
func (v *Validator) Validate(i interface{}) error {
	if err := v.Validator.Struct(i); err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)
		var errs []string
		for _, fieldError := range validationErrors {
			switch fieldError.Tag() {
			case "required":
				errs = append(errs, fmt.Sprintf("%s is required", fieldError.Field()))
			case "eq":
				errs = append(errs, fmt.Sprintf("%s must be equals to %s", fieldError.Field(), fieldError.Param()))
			default:
				errs = append(errs, err.Error())
			}
		}
		return echo.NewHTTPError(http.StatusBadRequest, errs)
	}

	return nil
}

// startHTTPServer create a new instance echo http server.
func startHTTPServer(rewardService v1.RewardService) {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "-" {
			return ""
		}
		return name
	})

	e.Validator = &Validator{Validator: validate}

	RegisterRoutes(e, rewardService)

	// Start server
	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	_ = e.Server.Shutdown(context.Background())
}

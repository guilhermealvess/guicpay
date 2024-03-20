package http

import (
	"log"
	"net/http"

	_ "github.com/guilhermealvess/guicpay/docs"
	"github.com/guilhermealvess/guicpay/internal/properties"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func NewServer(h *accountHandler) *echo.Echo {
	_, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}

	server := echo.New()
	server.Use(otelecho.Middleware("my-server"))
	server.GET("/docs/*", echoSwagger.WrapHandler)
	server.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "PONG\n")
	})

	api := server.Group("/api")
	api.POST("/accounts", h.CreateAccount)
	api.GET("/accounts", h.List, validateTokenMiddleware)
	api.GET("/accounts/me", h.Fetch, validateTokenMiddleware)

	api.POST("/transactions/deposit", h.AccountDeposit, validateTokenMiddleware)
	api.POST("/transactions/transfer", h.AccountTransfer, validateTokenMiddleware)

	api.POST("/auth", h.Auth)

	return server
}

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := zipkin.New(properties.Props.TraceCollectorURL)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

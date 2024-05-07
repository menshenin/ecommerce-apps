// Package app Приложение
package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"route256.ozon.ru/project/cart/internal/clients/loms"
	"route256.ozon.ru/project/cart/internal/clients/product"
	"route256.ozon.ru/project/cart/internal/httphandlers"
	lomspb "route256.ozon.ru/project/cart/internal/pkg/pb/loms"
	"route256.ozon.ru/project/cart/internal/repository/inmemorycart"
	"route256.ozon.ru/project/cart/internal/service/cart"
)

// App Приложение
type App struct {
	server *http.Server
}

// New Конструктор
func New(conf *Config) *App {
	logger := conf.Logger

	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(nil,
			otelhttp.WithTracerProvider(conf.TracerProvider),
			otelhttp.WithMeterProvider(conf.MeterProvider),
			otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				return fmt.Sprintf("%s %s", r.Method, r.URL.String())
			}),
			otelhttp.WithPropagators(
				propagation.NewCompositeTextMapPropagator(
					propagation.TraceContext{}, propagation.Baggage{})),
		),
	}
	productClient := product.New(
		conf.ProductServiceURL,
		product.WithHTTPClient(httpClient),
		product.WithToken(conf.ProductServiceToken),
		product.WithRetryMiddleware(conf.ProductServiceRetryCount),
		product.WithRPSLimit(conf.ProductServiceRPSLimit))
	cartRepo := inmemorycart.New()

	grpcConn, err := grpc.Dial(
		conf.LomsServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(
				otelgrpc.WithPropagators(propagation.NewCompositeTextMapPropagator(
					propagation.TraceContext{}, propagation.Baggage{},
				)),
				otelgrpc.WithMeterProvider(conf.MeterProvider),
				otelgrpc.WithTracerProvider(conf.TracerProvider))))

	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	lomsClient := loms.New(lomspb.NewLomsClient(grpcConn))

	cartService, err := cart.New(cartRepo, productClient, lomsClient)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	mux := http.NewServeMux()
	mux.Handle("POST /user/{user_id}/cart/{sku_id...}", httphandlers.AddItem(cartService))
	mux.Handle("DELETE /user/{user_id}/cart/{sku_id...}", httphandlers.DeleteItem(cartService))
	mux.Handle("DELETE /user/{user_id}/cart", httphandlers.Clear(cartService))
	mux.Handle("GET /user/{user_id}/cart", httphandlers.ListItems(cartService))
	mux.Handle("POST /cart/checkout", httphandlers.Checkout(cartService))

	mux.Handle("/metrics", promhttp.Handler())
	otelMW := otelhttp.NewMiddleware("serve",
		otelhttp.WithMeterProvider(conf.MeterProvider),
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		}),
		otelhttp.WithTracerProvider(conf.TracerProvider),
		otelhttp.WithFilter(func(request *http.Request) bool {
			return request.URL.Path != "/metrics"
		}),
	)

	server := &http.Server{
		Addr:              conf.ListenAddr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           otelMW(httphandlers.WithLog(logger)(mux.ServeHTTP)),
	}

	return &App{
		server: server,
	}
}

// Run Запуск приложения
func (a *App) Run() error {
	return a.server.ListenAndServe()
}

// Stop Остановка приложения
func (a *App) Stop(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

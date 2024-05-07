// Package app Приложение
package app

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	lomsServer "route256.ozon.ru/project/loms/internal/grpc/loms"
	"route256.ozon.ru/project/loms/internal/kafka/producer"
	lomspb "route256.ozon.ru/project/loms/internal/pkg/pb/loms"
	orderRepository "route256.ozon.ru/project/loms/internal/repository/dborder"
	stockRepository "route256.ozon.ru/project/loms/internal/repository/dbstock"
	orderService "route256.ozon.ru/project/loms/internal/service/order"
	stockService "route256.ozon.ru/project/loms/internal/service/stock"
)

// App Приложение
type App struct {
	config        Config
	grpcServer    *grpc.Server
	httpServer    *http.Server
	eventProducer *producer.Producer
}

// New Конструктор
func New(conf Config) (*App, error) {
	stockRepo := stockRepository.New(conf.MasterConnect, conf.SlaveConnect)
	orderRepo := orderRepository.New(conf.MasterConnect, conf.SlaveConnect)

	kafkaClient, err := kgo.NewClient(
		kgo.SeedBrokers(conf.KafkaAddr),
		kgo.DefaultProduceTopic(conf.KafkaTopic))
	if err != nil {
		return nil, err
	}
	eventProducer := producer.New(kafkaClient, slog.Default(), conf.TraceProvider)

	orders := orderService.New(orderRepo, stockRepo, eventProducer)
	stocks := stockService.New(stockRepo)
	lomsServ := lomsServer.New(orders, stocks)
	server := grpc.NewServer(conf.GrpcServerOpts...)
	lomspb.RegisterLomsServer(server, lomsServ)

	mux := runtime.NewServeMux()
	err = mux.HandlePath(http.MethodGet, "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		promhttp.Handler().ServeHTTP(w, r)
	})
	if err != nil {
		return nil, err
	}
	err = lomspb.RegisterLomsHandlerServer(context.Background(), mux, lomsServ)
	if err != nil {
		return nil, err
	}

	return &App{
		config:     conf,
		grpcServer: server,
		httpServer: &http.Server{
			Addr:              conf.HTTPGatewayAddr,
			Handler:           mux,
			ReadHeaderTimeout: 3 * time.Second,
		},
	}, nil
}

// Run Запуск приложения
func (a *App) Run() error {
	g := errgroup.Group{}
	g.Go(func() error {
		listener, err := net.Listen("tcp", a.config.ListenAddr)
		if err != nil {
			return err
		}
		return a.grpcServer.Serve(listener)
	})
	g.Go(func() error {
		return a.httpServer.ListenAndServe()
	})

	return g.Wait()
}

// Stop Остановка приложения
func (a *App) Stop(ctx context.Context) error {
	a.grpcServer.GracefulStop()
	a.eventProducer.Stop()

	return multierr.Combine(
		a.config.MasterConnect.Close(ctx),
		a.config.SlaveConnect.Close(ctx))
}

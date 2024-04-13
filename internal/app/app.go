package app

import (
	"context"
	"fmt"
	grpccli "github.com/Fruitfulfriends-REST-API-server/internal/clients/grpc"
	"github.com/Fruitfulfriends-REST-API-server/internal/config"
	"github.com/Fruitfulfriends-REST-API-server/internal/rest"
	"github.com/Fruitfulfriends-REST-API-server/internal/rest/handlers"
	"github.com/chatex-com/di-container"
	"github.com/chatex-com/process-manager"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

type Application struct {
	log       *log.Entry
	sigChan   <-chan os.Signal
	manager   *process.Manager
	container *di.Container
	cfg       *config.Config
	Done      chan struct{}
}

// NewApplication создает новый экземпляр приложения
func New(cfg *config.Config, log *log.Entry) (*Application, error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	a := &Application{
		log:       log,
		sigChan:   sigChan,
		cfg:       cfg,
		manager:   process.NewManager(),
		container: di.NewContainer(),
	}

	err := a.bootstrap()
	if err != nil {
		return nil, fmt.Errorf("failed to create new application: %w", err)
	}

	return a, nil
}

func (a *Application) Run() {
	a.Done = make(chan struct{}) // Инициализируем канал done
	a.manager.StartAll()
	go a.registerShutdown()
}

func (a *Application) bootstrap() error {
	if err := a.initGRPCWorker(); err != nil {
		return fmt.Errorf("failed to init GRPC worker: %w", err)
	}

	if err := a.initRestWorker(); err != nil {
		return fmt.Errorf("failed to init rest worker: %w", err)
	}

	return nil
}

func (a *Application) initGRPCWorker() error {
	const op = "Application.initAdminGRPC"
	a.log.WithField("operation", op).Info("initializing admin grpc")

	apiService, err := grpccli.New(context.Background(), a.log, a.cfg.Clients.GRPC.Address, a.cfg.Clients.GRPC.Timeout, a.cfg.Clients.GRPC.RetriesCount)
	if err != nil {
		return fmt.Errorf("failed to create grpc client: %w", err)
	}

	a.container.Set(apiService)
	return nil
}

func (a *Application) initRestWorker() error {
	const op = "Application.initRestWorker"
	a.log.WithField("operation", op).Info(("initializing rest worker"))

	var apiService *grpccli.Client
	if err := a.container.Load(&apiService); err != nil {
		return fmt.Errorf("%s: failed to load grpc client: %w", op, err)
	}

	apiHandlers := []handlers.APIHandler{
		handlers.NewAuthHandler(apiService, a.log.Logger, a.cfg.AppID),
	}

	w := rest.NewWorker(
		&a.cfg.HTTPServer,
		a.log.Logger,
		apiHandlers,
	)

	cb := process.NewCallbackWorker("Rest server", w.Start)
	a.manager.AddWorker(cb)

	return nil
}

func (a *Application) registerShutdown() {
	const op = "Application.registerShutdown"

	go func(manager *process.Manager) {
		<-a.sigChan
		manager.StopAll()
		close(a.Done) // Закрываем канал done после остановки всех воркеров
		a.log.WithField("operation", op).Info("registering shutdown")
	}(a.manager)

	a.manager.AwaitAll()
}

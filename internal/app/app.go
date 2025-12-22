package app

import (
	"context"
	"log"

	"event-driven-architecture/internal/adapter/appctx"
	"event-driven-architecture/internal/adapter/db/postgres"
	"event-driven-architecture/internal/adapter/db/session"
	"event-driven-architecture/internal/adapter/event/broker"
	"event-driven-architecture/internal/adapter/event/publisher"
	"event-driven-architecture/internal/adapter/hasher"
	"event-driven-architecture/internal/adapter/id"
	"event-driven-architecture/internal/adapter/logger"
	"event-driven-architecture/internal/adapter/server"
	"event-driven-architecture/internal/infrastructure"
	"event-driven-architecture/internal/usecase"
	"event-driven-architecture/pkg/envloader"
)

type App struct {
	log            usecase.Logger
	sessionManager usecase.SessionManager
	appCtxManager  usecase.AppCtxManager
	sseBroker      usecase.EventBroker

	createUserUseCase *usecase.CreateUserUseCase
}

func NewApp() *App {
	// infrastructure
	natsConn, err := infrastructure.NewNATS(envloader.MustGetString("NATS_URL"))
	if err != nil {
		log.Fatalln("не удалось подключиться к NATS:", err)
	}

	postgresConn, err := infrastructure.NewPostgres(envloader.MustGetString("POSTGRES_DNS"))
	if err != nil {
		log.Fatalln("не удалось подключиться к Postgres:", err)
	}

	// adapters
	log := logger.NewSlogLogger()
	sessionManager := session.NewSessionManager(postgresConn)
	appCtxManager := appctx.NewAppCtxManager(envloader.GetDuration("APP_CTX_TIMEOUT", usecase.AppCtxDefaultTimeout))

	sseBroker := broker.NewSSEBroker(log)

	natsPublisher := publisher.NewNATSPublisher(natsConn)
	ssePublisher := publisher.NewSSEPublisher(sseBroker)

	uuidGenerator := id.NewUUIDGenerator()
	bcryptHasher := hasher.NewBCrypt(envloader.GetInt("BCRYPT_COST", 12))

	userRepo := postgres.NewUser()

	// use-cases
	createUserUseCase := usecase.NewCreateUser(
		log,
		sessionManager,
		appCtxManager,
		publisher.NewMultiPublisher(natsPublisher, ssePublisher),
		uuidGenerator,
		bcryptHasher,
		userRepo,
	)

	return &App{
		log:            log,
		sessionManager: sessionManager,
		appCtxManager:  appCtxManager,
		sseBroker:      sseBroker,

		createUserUseCase: createUserUseCase,
	}
}

func (a *App) Run() {
	defer a.sessionManager.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: добавить обработку CTRL+C

	a.server().Run(ctx)
}

func (a *App) server() *server.HTTPServer {
	return server.NewHTTPServer(
		a.log,
		envloader.MustGetInt("PORT"),
		a.sseBroker,
		a.createUserUseCase,
	)
}

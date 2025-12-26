package app

import (
	broker "event-driven-architecture/internal/adapter/event/broker/sse_broker"
	"event-driven-architecture/internal/adapter/event/publisher"
	"event-driven-architecture/internal/adapter/hasher"
	"event-driven-architecture/internal/adapter/id"
	"event-driven-architecture/internal/adapter/logger"
	"event-driven-architecture/internal/adapter/repo/multi"
	"event-driven-architecture/internal/adapter/repo/postgres"
	"event-driven-architecture/internal/adapter/repo/rediscache"
	"event-driven-architecture/internal/app/context/appctx"
	"event-driven-architecture/internal/app/transaction"
	"event-driven-architecture/internal/infrastructure"
	"event-driven-architecture/internal/usecase"
)

type App struct {
	Log           usecase.Logger
	TxManager     usecase.TransactionManager
	AppCtxManager usecase.AppCtxManager
	SSEBroker     usecase.EventBroker

	AuthenticateSessionUseCase *usecase.AuthenticateSessionUseCase
	LoginUserUseCase           *usecase.LoginUserUseCase
	CreateUserUseCase          *usecase.CreateUserUseCase
}

func Build(cfg Config) (*App, error) {
	// --- infrastructure ---
	log := logger.NewSlogLogger()

	pgConn, err := infrastructure.NewPostgres(cfg.Postgres.DNS)
	if err != nil {
		return nil, err
	}

	redisClient, err := infrastructure.NewRedis(cfg.Redis.Addr)
	if err != nil {
		return nil, err
	}

	// --- cross-cutting ---
	txManager := transaction.NewManager(pgConn)
	appCtxManager := appctx.NewManager(cfg.App.ContextTimeout)
	uuidGen := id.NewUUIDGenerator()
	bcryptHasher := hasher.NewBCrypt(cfg.Security.BcryptCost)

	// --- repositories ---
	// postgres
	pgUserRepo := postgres.NewUser()

	// redis
	redisSessionIndexRepo := rediscache.NewSessionIndex(redisClient)

	// multi
	multiSessionRepo := multi.NewMultiSession(
		rediscache.NewSession(redisClient),
		postgres.NewSession(),
	)

	// -- events ---
	sseBroker := broker.NewSSEBroker(log)
	ssePublisher := publisher.NewSSEPublisher(sseBroker, redisSessionIndexRepo)

	// --- use-cases ---
	createUserUseCase := usecase.NewCreateUser(
		log,
		ssePublisher,
		uuidGen,
		bcryptHasher,
		pgUserRepo,
	)

	loginUserUseCase := usecase.NewLoginUser(
		log,
		uuidGen,
		bcryptHasher,
		pgUserRepo,
		multiSessionRepo,
		redisSessionIndexRepo,
	)

	authSessionUseCase := usecase.NewAuthenticateSession(
		log,
		multiSessionRepo,
	)

	return &App{
		Log:           log,
		TxManager:     txManager,
		AppCtxManager: appCtxManager,
		SSEBroker:     sseBroker,

		AuthenticateSessionUseCase: authSessionUseCase,
		LoginUserUseCase:           loginUserUseCase,
		CreateUserUseCase:          createUserUseCase,
	}, nil
}

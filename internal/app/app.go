package app

import (
	"context"
	"log"

	"event-driven-architecture/internal/adapter/appctx"
	broker "event-driven-architecture/internal/adapter/event/broker/sse_broker"
	"event-driven-architecture/internal/adapter/event/publisher"
	"event-driven-architecture/internal/adapter/hasher"
	"event-driven-architecture/internal/adapter/id"
	"event-driven-architecture/internal/adapter/logger"
	"event-driven-architecture/internal/adapter/repo/multi"
	"event-driven-architecture/internal/adapter/repo/postgres"
	"event-driven-architecture/internal/adapter/repo/rediscache"
	"event-driven-architecture/internal/adapter/server"
	"event-driven-architecture/internal/adapter/transaction"
	"event-driven-architecture/internal/infrastructure"
	"event-driven-architecture/internal/usecase"
	"event-driven-architecture/pkg/envloader"
)

type App struct {
	log                usecase.Logger
	transactionManager usecase.TransactionManager
	appCtxManager      usecase.AppCtxManager
	sseBroker          usecase.EventBroker

	authenticateSessionUseCase *usecase.AuthenticateSessionUseCase
	loginUserUseCase           *usecase.LoginUserUseCase
	createUserUseCase          *usecase.CreateUserUseCase
}

func NewApp() *App {
	// infrastructure
	postgresConn, err := infrastructure.NewPostgres(envloader.MustGetString("POSTGRES_DNS"))
	if err != nil {
		log.Fatalln("не удалось подключиться к Postgres:", err)
	}

	redisClient, err := infrastructure.NewRedis(envloader.MustGetString("REDIS_ADDR"))
	if err != nil {
		log.Fatalln("не удалось подключиться к Redis:", err)
	}

	// adapters
	log := logger.NewSlogLogger()
	transactionManager := transaction.NewTransactionManager(postgresConn)
	appCtxManager := appctx.NewAppCtxManager(envloader.GetDuration("APP_CTX_TIMEOUT", usecase.AppCtxDefaultTimeout))

	uuidGenerator := id.NewUUIDGenerator()
	bcryptHasher := hasher.NewBCrypt(envloader.GetInt("BCRYPT_COST", 12))

	userPostgresRepo := postgres.NewUser()
	sessionPostgresRepo := postgres.NewSession()
	sessionRedisRepo := rediscache.NewSession(redisClient)
	sessionIndexRedisRepo := rediscache.NewSessionIndex(redisClient)

	sseBroker := broker.NewSSEBroker(log)
	ssePublisher := publisher.NewSSEPublisher(sseBroker, sessionIndexRedisRepo)

	// use-cases
	authenticateSessionUseCase := usecase.NewAuthenticateSession(
		log,
		sessionRedisRepo,
	)

	loginUserUseCase := usecase.NewLoginUser(
		log,
		uuidGenerator,
		bcryptHasher,
		userPostgresRepo,
		multi.NewMultiSession(sessionRedisRepo, sessionPostgresRepo),
		sessionIndexRedisRepo,
	)

	createUserUseCase := usecase.NewCreateUser(
		log,
		ssePublisher,
		uuidGenerator,
		bcryptHasher,
		userPostgresRepo,
	)

	return &App{
		log:                log,
		transactionManager: transactionManager,
		appCtxManager:      appCtxManager,
		sseBroker:          sseBroker,

		authenticateSessionUseCase: authenticateSessionUseCase,
		loginUserUseCase:           loginUserUseCase,
		createUserUseCase:          createUserUseCase,
	}
}

func (a *App) Run() {
	defer a.transactionManager.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: добавить обработку CTRL+C

	a.server().Run(ctx)
}

func (a *App) server() *server.HTTPServer {
	return server.NewHTTPServer(
		a.log,
		envloader.MustGetInt("PORT"),
		a.transactionManager,
		a.appCtxManager,
		a.sseBroker,
		a.authenticateSessionUseCase,
		a.loginUserUseCase,
		a.createUserUseCase,
	)
}

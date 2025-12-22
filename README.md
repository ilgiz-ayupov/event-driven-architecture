## Общая структура

```
/cmd
    main.go

/internal
  /app
    app.go

  /domain
    user.go
    events.go

  /usecase
    create_user.go
    ports.go

  /adapter
    /httpserver
      user_handler.go
    /event
      nats_publisher.go
      nats_subscriber.go
    /logger

  /infrastructure
    nats.go

/pkg
```

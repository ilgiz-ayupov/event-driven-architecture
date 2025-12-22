# Принципы Чистой Архитектуры и их применение

Чистая Архитектура (Clean Architecture), популяризированная Робертом С. Мартином (Uncle Bob), предлагает подход к организации кода, который делает программные системы независимыми от фреймворков, баз данных и внешних агентов. Её основная цель — создать систему, которая легко тестируется, поддерживается и развивается, сосредотачиваясь на бизнес-правилах приложения, а не на технических деталях их реализации.

## Основные Принципы Чистой Архитектуры

Чистая Архитектура основывается на нескольких ключевых принципах, которые обеспечивают разделение ответственности и независимость слоев. Центральным элементом является Принцип Инверсии Зависимостей (Dependency Inversion Principle, DIP), который гласит, что высокоуровневые модули не должны зависеть от низкоуровневых модулей; оба должны зависеть от абстракций. Абстракции не должны зависеть от деталей; детали должны зависеть от абстракций.

### Независимость от Фреймворков (Framework Independence)

Система должна быть независима от использования конкретных фреймворков. Это означает, что бизнес-логика приложения не должна "знать" о том, какой веб-фреймворк (например, Gin), какая база данных или какая система очередей сообщений используется. Если фреймворк изменится или понадобится использовать другой, это не должно требовать переписывания основной логики приложения.

*Пример 1 - Веб-сервер:*

Представьте, что у вас есть бизнес-логика для обработки заказа. В традиционном подходе эта логика может быть тесно связана с HTTP-хендлерами Gin, например, напрямую получая параметры из c.Param("id") и записывая ответ через c.JSON(). В Чистой Архитектуре бизнес-логика (например, OrderService.ProcessOrder) принимает только чистые данные (структуры Go) и возвращает чистые данные. Адаптер Gin (инфраструктурный слой) будет отвечать за преобразование HTTP-запроса в структуру данных для OrderService и преобразование результата OrderService в HTTP-ответ. Если вы решите перейти с Gin на Echo или Fasthttp, вам нужно будет изменить только адаптеры, а бизнес-логика останется нетронутой.

*Пример 2 - Уведомления:*

Допустим, у вас есть сервис уведомлений, который отправляет сообщения пользователям. В нечистой архитектуре, этот сервис мог бы напрямую вызывать методы сторонней библиотеки для отправки SMS или Email. В Чистой Архитектуре, ваш бизнес-сервис уведомлений будет зависеть от абстракции Notifier (интерфейса), который определяет метод Send(message string, recipient string) error. Конкретная реализация SMSService или EmailService будет реализовать этот интерфейс и будет находиться в инфраструктурном слое. Бизнес-логика будет работать только с интерфейсом, не зная, как именно сообщение будет отправлено.

### Тестируемость (Testability)

Бизнес-правила могут быть протестированы без пользовательского интерфейса, базы данных, сервера или любого другого внешнего элемента. Это значительно упрощает процесс разработки и обеспечивает высокую надежность кода.

*Пример 1 - Расчёт скидок:*

У вас есть логика расчёта скидок для продукта. В Чистой Архитектуре эта логика инкапсулирована в доменном или прикладном слое. Вы можете написать юнит-тесты, которые напрямую вызывают функцию CalculateDiscount(productPrice float64, userCategory string) с различными входными данными, не беспокоясь о том, как эти данные будут получены из базы данных или переданы через HTTP-запрос. Мокировать базы данных или HTTP-запросы не требуется для тестирования бизнес-логики.

*Пример 2 - Валидация данных:*

Предположим, у вас есть правила валидации для создания нового пользователя (например, длина пароля, формат email). В Чистой Архитектуре эта валидация находится в Use Case (прикладном слое) или даже в доменной сущности. Вы можете протестировать эти правила, передавая различные структуры User напрямую в функцию валидации, без необходимости запускать весь веб-сервер или взаимодействовать с базой данных.

### Независимость от UI (UI Independence)

Изменение пользовательского интерфейса (например, переход с веб-приложения на мобильное приложение или API) не должно требовать изменения основной логики приложения.

*Пример 1 - API и CLI:*

У вас есть набор функций, управляющих задачами (создание, чтение, обновление, удаление). В Чистой Архитектуре эта логика находится в слоях Domain и Application. Вы можете создать адаптер для Gin, который предоставляет RESTful API для этих задач, а также адаптер для Cobra CLI, который позволяет управлять задачами из командной строки. При этом, основная логика работы с задачами остаётся неизменной.

*Пример 2 - Различные форматы вывода:*

Бизнес-логика генерирует некий отчёт в виде структуры данных. Адаптеры инфраструктурного слоя могут взять эту структуру и преобразовать её в JSON для веб-API, в CSV для экспорта, или отформатировать в HTML для веб-страницы. Логика генерации отчёта не зависит от формата его представления.

### Независимость от Баз Данных (Database Independence)

Бизнес-правила не должны "знать", какая база данных используется (SQL, NoSQL, ORM и т.д.). Переход на другую базу данных не должен влиять на доменный и прикладной слои.

*Пример 1 - Хранение пользователей:*

Ваш сервис управления пользователями использует интерфейс UserRepository, который определяет методы `GetUserByID(id string) (*User, error)` и `SaveUser(user*User) error`. Конкретная реализация PostgreSQLUserRepository будет использовать pgx или database/sql для взаимодействия с PostgreSQL, а MongoDBUserRepository будет использовать mongo-driver для MongoDB. Бизнес-логика (например, UserService) будет зависеть только от интерфейса UserRepository, не заботясь о деталях хранения.

*Пример 2 - Логирование событий:*

Сервис логирования событий в приложении использует интерфейс EventLogger, который определяет метод Log(event Event) error. Одна реализация может записывать события в файлы, другая — в Elasticsearch, третья — в облачное хранилище. Сервис, генерирующий события, вызывает метод Log интерфейса, не зная, куда именно будут записаны события.

### Независимость от Внешних Агентов (External Agency Independence)

Любые внешние сервисы, такие как облачные провайдеры, системы сообщений или сторонние API, также должны быть абстрагированы, чтобы их изменение не влияло на основные бизнес-правила.

*Пример 1 - Платежные шлюзы:*

Ваш сервис электронной коммерции обрабатывает платежи. Вместо того чтобы бизнес-логика напрямую вызывала Stripe API, она будет зависеть от интерфейса PaymentGateway, который определяет метод ProcessPayment(amount float64, currency string) (string, error). Реализации StripeGateway и PayPalGateway будут находиться в инфраструктурном слое и инкапсулировать детали взаимодействия с конкретными платежными системами.

*Пример 2 - Отправка уведомлений через различные каналы:*

Сервис может отправлять уведомления как по SMS через Twilio, так и по электронной почте через SendGrid. Бизнес-логика будет использовать общий интерфейс NotificationSender, а конкретные реализации для Twilio и SendGrid будут находиться в слое инфраструктуры, абстрагируя бизнес-логику от внешних деталей.

## Слои Чистой Архитектуры

Чистая Архитектура традиционно описывается как набор концентрических кругов, где каждый внешний круг зависит от внутренних, но внутренние круги ничего не знают о внешних. Это "Правило Зависимостей" (Dependency Rule).

### Сущности (Entities / Domain Layer)

Это самый внутренний слой. Он содержит бизнес-объекты и правила, которые являются наиболее общими и высокоуровневыми. Сущности инкапсулируют корпоративные бизнес-правила и являются ядром приложения. Они могут быть структурами данных с методами, которые содержат бизнес-логику, которая может использоваться многими приложениями. Они не должны зависеть от чего-либо извне.

*Пример - Пользователь (User):*

```go
package domain

import "errors"

// User представляет сущность пользователя в домене
type User struct {
 ID       string
 Email    string
 Password string // Хэшированный пароль
 IsActive bool
 Role     string // Например, "admin", "customer"
}

// NewUser создает нового пользователя с базовой валидацией
func NewUser(id, email, password string) (*User, error) {
 if id == "" || email == "" || password == "" {
  return nil, errors.New("ID, email и пароль не могут быть пустыми")
 }
 // Дополнительные правила домена, например, проверка формата email
 return &User{
  ID:       id,
  Email:    email,
  Password: password, // В реальном приложении здесь будет хэширование
  IsActive: true,
  Role:     "customer",
 }, nil
}

// UpdateEmail обновляет адрес электронной почты пользователя
// Это доменное правило, которое может включать валидацию
func (u *User) UpdateEmail(newEmail string) error {
 if newEmail == "" {
  return errors.New("email не может быть пустым")
 }
 // Можно добавить более сложную валидацию email
 u.Email = newEmail
 return nil
}
```

Здесь User — это сущность. NewUser и UpdateEmail содержат бизнес-правила, связанные с пользователем. Они не зависят от базы данных, веб-фреймворка или чего-либо еще.

### Варианты Использования (Use Cases / Application Layer)

Этот слой содержит специфичные для приложения бизнес-правила. Он оркестрирует поток данных к и от сущностей. Варианты использования зависят от сущностей, но сущности не зависят от вариантов использования. Этот слой управляет тем, как данные текут через систему, и делегирует работу сущностям. Он определяет интерфейсы для внешних элементов (баз данных, API и т.д.), которые затем реализуются в инфраструктурном слое.

*Пример - Создание Пользователя (CreateUserUseCase):*

```go
package application

import (
 "context"
 "errors"
 "fmt"
 "your_project/domain" // Предполагается, что domain пакет находится на том же уровне
)

// UserRepository определяет интерфейс для работы с хранилищем пользователей
// Это абстракция, от которой зависит Use Case
type UserRepository interface {
 Save(ctx context.Context, user *domain.User) error
 GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

// UserOutputPort определяет интерфейс для представления результатов
// Это абстракция, от которой зависит Use Case, чтобы уведомить о результате
type UserOutputPort interface {
 PresentUser(user *domain.User)
 PresentError(err error)
}

// CreateUserRequest содержит входные данные для создания пользователя
type CreateUserRequest struct {
 ID       string
 Email    string
 Password string
}

// CreateUserUseCase обрабатывает бизнес-логику создания нового пользователя
type CreateUserUseCase struct {
 userRepo   UserRepository
 outputPort UserOutputPort
}

// NewCreateUserUseCase создает новый экземпляр CreateUserUseCase
func NewCreateUserUseCase(repo UserRepository, op UserOutputPort) *CreateUserUseCase {
 return &CreateUserUseCase{
  userRepo:   repo,
  outputPort: op,
 }
}

// Execute выполняет логику создания пользователя
func (uc *CreateUserUseCase) Execute(ctx context.Context, req CreateUserRequest) {
 // Проверка существования пользователя с таким email
 existingUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
 if err != nil && err.Error() != "user not found" { // Предполагаем ошибку "user not found"
  uc.outputPort.PresentError(fmt.Errorf("ошибка при проверке существующего пользователя: %w", err))
  return
 }
 if existingUser != nil {
  uc.outputPort.PresentError(errors.New("пользователь с таким email уже существует"))
  return
 }

 // Создание сущности User
 user, err := domain.NewUser(req.ID, req.Email, req.Password)
 if err != nil {
  uc.outputPort.PresentError(fmt.Errorf("ошибка создания сущности пользователя: %w", err))
  return
 }

 // Сохранение пользователя через репозиторий
 if err := uc.userRepo.Save(ctx, user); err != nil {
  uc.outputPort.PresentError(fmt.Errorf("ошибка сохранения пользователя: %w", err))
  return
 }

 uc.outputPort.PresentUser(user)
}
```

Здесь CreateUserUseCase является вариантом использования. Он содержит бизнес-логику для создания пользователя, используя интерфейс UserRepository (который будет реализован в инфраструктурном слое) и интерфейс UserOutputPort для уведомления о результате. Обратите внимание, что application слой зависит от domain, но domain ничего не знает об application.

### Адаптеры Интерфейсов (Interface Adapters)

Этот слой содержит адаптеры, которые преобразуют данные из формата, наиболее удобного для внешнего мира (UI, базы данных, внешние сервисы), в формат, наиболее удобный для вариантов использования и сущностей, и наоборот. К ним относятся контроллеры, шлюзы (gateways) и презентеры.

*Пример - HTTP-контроллер (Gin Adapter) и Презентер:*

```go
package infrastructure

import (
 "context"
 "net/http"
 "your_project/application"
 "your_project/domain"

 "github.com/gin-gonic/gin"
)

// UserPresenter implements application.UserOutputPort
type UserPresenter struct {
 Ctx *gin.Context
}

// PresentUser форматирует успешный результат для HTTP-ответа
func (p *UserPresenter) PresentUser(user *domain.User) {
 p.Ctx.JSON(http.StatusCreated, gin.H{
  "id":      user.ID,
  "email":   user.Email,
  "message": "User created successfully",
 })
}

// PresentError форматирует ошибку для HTTP-ответа
func (p *UserPresenter) PresentError(err error) {
 statusCode := http.StatusInternalServerError
 if err.Error() == "пользователь с таким email уже существует" {
  statusCode = http.StatusConflict
 } else if err.Error() == "ID, email и пароль не могут быть пустыми" || err.Error() == "email не может быть пустым" {
  statusCode = http.StatusBadRequest
 }
 p.Ctx.JSON(statusCode, gin.H{"error": err.Error()})
}

// GinUserController является адаптером для Gin
type GinUserController struct {
 createUserUseCase *application.CreateUserUseCase
}

// NewGinUserController создает новый экземпляр GinUserController
func NewGinUserController(createUserUC *application.CreateUserUseCase) *GinUserController {
 return &GinUserController{
  createUserUseCase: createUserUC,
 }
}

// CreateUserHandler обрабатывает HTTP-запрос для создания пользователя
func (ctrl *GinUserController) CreateUserHandler(c *gin.Context) {
 var req application.CreateUserRequest
 if err := c.ShouldBindJSON(&req); err != nil {
  c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
  return
 }

 // Создаем презентер для данного запроса
 presenter := &UserPresenter{Ctx: c}

 // Вызываем Use Case, передавая ему контекст и презентер
 ctrl.createUserUseCase.Execute(c.Request.Context(), req)
 // Важно: Use Case сам вызовет методы Presenter для ответа,
 // поэтому здесь не нужно c.JSON().
 // Однако, для простоты примера, мы оставим Presenter
 // как способ уведомления Gin-хендлера, а не прямого ответа.
 // В более сложной системе Presenter может быть использован для подготовки ViewModels.
 // Здесь мы его используем для прямого ответа через Ctx.JSON, что допустимо
 // для простых API, но стоит помнить о более чистых паттернах.
}

// UserRepositoryImpl implements application.UserRepository
type UserRepositoryImpl struct {
 // Здесь будут детали для работы с базой данных, например, *sql.DB
 // Для примера, используем простой map
 users map[string]*domain.User
}

func NewUserRepositoryImpl() *UserRepositoryImpl {
 return &UserRepositoryImpl{
  users: make(map[string]*domain.User),
 }
}

func (r *UserRepositoryImpl) Save(ctx context.Context, user *domain.User) error {
 // Имитация сохранения в базу данных
 r.users[user.ID] = user
 fmt.Printf("User %s saved to 'database'\n", user.ID)
 return nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
 // Имитация получения из базы данных
 for _, user := range r.users {
  if user.Email == email {
   return user, nil
  }
 }
 return nil, errors.New("user not found") // Предполагаемая ошибка "user not found"
}
```

Здесь GinUserController является адаптером, который принимает HTTP-запросы, преобразует их в формат CreateUserRequest и передает в CreateUserUseCase. UserPresenter преобразует результаты работы CreateUserUseCase в HTTP-ответы. UserRepositoryImpl реализует интерфейс UserRepository из слоя application, предоставляя конкретную реализацию для работы с хранилищем (в данном случае, имитация через map). Обратите внимание, что infrastructure слой зависит от application и domain.

### Фреймворки и Драйверы (Frameworks & Drivers / Infrastructure Layer)

Это самый внешний слой, содержащий конкретные реализации для баз данных, веб-фреймворков (Gin), UI и других внешних устройств. Этот слой не должен содержать бизнес-логику и должен быть легко заменяемым. Он зависит от всех внутренних слоев.

*Пример - Точка входа приложения (main.go):*

```go
package main

import (
 "log"
 "your_project/application"
 "your_project/infrastructure"

 "github.com/gin-gonic/gin"
)

func main() {
 r := gin.Default()

 // Инициализация репозитория (конкретная реализация)
 userRepo := infrastructure.NewUserRepositoryImpl()

 // Инициализация Use Case
 createUserUseCase := application.NewCreateUserUseCase(userRepo, &infrastructure.UserPresenter{}) // Презентер будет создан в хендлере

 // Инициализация контроллера
 userController := infrastructure.NewGinUserController(createUserUseCase)

 // Настройка маршрутов
 r.POST("/users", userController.CreateUserHandler)

 log.Println("Сервер запущен на :8080")
 if err := r.Run(":8080"); err != nil {
  log.Fatalf("Ошибка запуска сервера: %v", err)
 }
}
```

В main.go мы "собираем" все части приложения. Здесь создаются конкретные реализации репозиториев и контроллеров, и они связываются вместе. Это точка, где внешние зависимости (Gin, map в качестве имитации БД) инициализируются и подключаются к внутренним слоям.

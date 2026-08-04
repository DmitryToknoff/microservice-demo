## Архитектура системы
                               
                               ┌──────────────────────┐
                               │     Client / Web     │
                               └──────────┬───────────┘
                                          │ HTTP / REST (JSON)
                                          v
                               ┌──────────────────────┐
                               │     API Gateway      │
                               └──────────┬───────────┘
                                          │
                   ┌──────────────────────┴──────────────────────┐
                   │ gRPC (Protobuf)                             │ gRPC (Protobuf)
                   v                                             v
        ┌──────────────────────┐                      ┌──────────────────────┐
        │    Order Service     │                      │   Delivery Service   │
        └──────────┬───────────┘                      └──────────┬───────────┘
                   │                                             ▲
                   │ Producer                                    │ Consumer (At-Least-Once)
                   │                                             │
                   └─────────────────► Kafka ◄───────────────────┘
                                 (Topic: order.created)



### Стек технологий
```text
⚬	Язык программирования: Go 1.24+
⚬	Транспорт & Протоколы: gRPC, Protobuf, REST (HTTP/JSON), segmentio/kafka-go
⚬	Базы данных & Кэш: PostgreSQL (pgx/v5 pool), Redis
⚬	Брокер сообщений: Apache Kafka (Кластер из 3-х брокеров) + Zookeeper
⚬	Конфигурирование: kelseyhightower/envconfig
⚬	Логирование: uber-go/zap
⚬	Контейнеризация: Docker, Docker Compose (Multi-stage builds)
```

## Структура проекта

```text
.
├── cmd/                        # Точки входа в приложения
│   ├── gateway/                # API Gateway
│   ├── order/                  # Order Service
│   └── delivery/               # Delivery Service
│
├── internal/                   # Внутренняя бизнес-логика (Clean Architecture)
│   ├── gateway/                # HTTP Handlers, gRPC Clients
│   ├── order/                  # Domain, Repository, Service, gRPC Transport
│   └── delivery/               # Domain, Repository, Service, Kafka Consumer
│
├── pkg/                        # Общие переиспользуемые модули
│   ├── logger/                 # Логгер на базе Zap
│   ├── postgres/               # Пул подключений PostgreSQL
│   ├── redis/                  # Клиент Redis
│   └── pb/                     # Сгенерированный gRPC / Protobuf код
│
├── migrations/                 # SQL-миграции баз данных
├── api/proto/                  # Protobuf схемы (.proto)
├── docker-compose.yml          # Оркестрация всей инфраструктуры
└── go.mod                      # Зависимости Go
```

## Быстрый запуск

1.	Клонируйте репозиторий:
```bash
  git clone https://github.com/DmitryToknoff/microservice-demo.git
```
```bash
  cd microservice-demo
```

2.	Запустите сборку и старт всех сервисов в Docker:

```bash
docker compose up -d --build
```

## Примеры использования (REST API)

1. Создание нового заказа
При создании заказа Order Service сохранит запись в БД и вызовет событие в Kafka. Delivery Service автоматически вычитает это событие и создаст запись о доставке.

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 101,
    "amount": 1499.90
  }'
```

Ответ (201 Created / 200 OK):

```json
{
  "id": 1,
  "user_id": 101,
  "amount": 1499.9,
  "status": "CREATED"
}
```

2. Получение информации о заказе
Данный запрос проверяет кэш Redis. При отсутствии данных запрос идёт в PostgreSQL, а полученное значение кэшируется на заданный TTL.

```bash
curl -X GET http://localhost:8080/api/v1/orders/1 \
  -H "Accept: application/json"
```

Ответ (200 OK):
```json
{
  "id": 1,
  "user_id": 101,
  "amount": 1499.9,
  "status": "CREATED"
}
```

3. Получение статуса доставки
Возвращает текущую информацию о доставке, асинхронно созданной при вычитывании Kafka-события.

```bash
curl -X GET http://localhost:8080/api/v1/deliveries/1 \
  -H "Accept: application/json"
```

Ответ (200 OK):

```json
{
  "order_id": 1,
  "status": "PROCESSING",
  "address": "Default delivery address"
} 
```

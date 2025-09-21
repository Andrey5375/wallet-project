# Wallet Service

Приложение на Go для работы с кошельками через REST API. Позволяет делать депозиты, снимать средства и получать баланс кошелька. Работает с PostgreSQL и запускается через Docker Compose.  

## Стек технологий
- **Golang** (чистая архитектура, конкурентные операции, тесты)  
- **PostgreSQL** (хранение балансов, транзакционная обработка)  
- **Docker & Docker Compose** (изолированный запуск сервиса и БД)  
- **REST API** (простое взаимодействие с клиентами)  
- **hey** (нагрузочное тестирование, проверка 1000+ RPS)  

## Архитектура
Проект разделён на слои:  
- `internal/app` — сервер и маршруты  
- `internal/service` — бизнес-логика (депозиты, снятия, проверки)  
- `internal/repo` — работа с БД (PostgreSQL)  
- `pkg/config` — загрузка переменных окружения из `config.env`  
- `service_test` — модульные и конкурентные тесты  

## Структура проекта
```bash
wallet-project/
│── cmd/
│   └── server/                # Точка входа в приложение (main.go)
│
│── internal/
│   ├── api/                   # REST API (роутинг, эндпоинты, v1/wallet)
│   ├── app/                   # Инициализация HTTP-сервера и зависимостей
│   ├── service/               # Бизнес-логика (операции с кошельками: deposit, withdraw, balance)
│   ├── repo/                  # Работа с PostgreSQL (хранение и транзакции по кошелькам)
│
│── pkg/
│   ├── config/                # Загрузка конфигурации из config.env
│   └── logger/                # Логирование (уровни логов, формат вывода)
│
│── service_test/              # Тесты (юнит + конкурентные тесты на депозит/снятие)
│
│── migrations/                # SQL-скрипты для создания и миграции таблиц
│
│── Dockerfile                 # Docker-образ приложения
│── docker-compose.yml         # Поднятие сервиса и PostgreSQL
│── config.env                 # Переменные окружения (DB_URL, PORT и т.д.)
│── README.md                  # Документация по запуску проекта

```

## Запуск проекта
1. Склонировать репозиторий:  
```bash
git clone https://github.com/Andrey5375/wallet-project.git
cd wallet-project
```
2. Создать файл config.env с переменными окружения:
```bash
PORT=8080
DATABASE_URL=postgres://user:password@db:5432/wallet?sslmode=disable
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=wallet
```

3. Поднять сервис и базу через Docker Compose:
```bash
docker-compose up --build
```

4. Приложение будет доступно на:
```bash
http://localhost:8080
```

API
POST /api/v1/wallet

Выполняет операцию DEPOSIT или WITHDRAW.
```bash
Request:
{
  "walletId": "11111111-1111-1111-1111-111111111111",
  "operationType": "DEPOSIT",
  "amount": 1000
}


Response (успех):
{
  "status": "ok"
}

Response (ошибка, недостаточно средств):
{
  "error": "cannot withdraw: insufficient funds"
}
```

GET /api/v1/wallets/{walletId}

Возвращает текущий баланс кошелька.

```bash
Response:

{
  "balance": 12345
}
```

Тесты

Конкурентные тесты находятся в папке service_test.
Покрывают:

Одновременные депозиты

Одновременные снятия

Обработку ошибок при недостатке средств

Запуск тестов:
```bash
go test -v ./service_test
```

Нагрузочное тестирование

Пример проверки 1000 запросов с 50 конкурентными потоками:
```bash
hey -n 1000 -c 50 -m POST -H "Content-Type: application/json" -d '{"walletId":"11111111-1111-1111-1111-111111111111","operationType":"DEPOSIT","amount":10}' http://localhost:8080/api/v1/wallet
```
Результат: 0 ошибок (50X нет), ~2000 RPS на MacBook Pro.

Миграции

При старте контейнера автоматически создаётся таблица wallets из папки migrations. Можно расширить новыми SQL-скриптами для дополнительного функционала.

Особенности

Обработка 1000+ RPS по одному кошельку без ошибок, все операции атомарные (используются транзакции PostgreSQL), переменные окружения загружаются из config.env, полное покрытие тестами для конкурентной среды, возможность масштабирования через docker-compose scale app=N.


Примеры команд cURL

Пополнить кошелёк:
```bash
curl -X POST http://localhost:8080/api/v1/wallet
 -H "Content-Type: application/json" -d '{"walletId":"11111111-1111-1111-1111-111111111111","operationType":"DEPOSIT","amount":500}'

Снять деньги:
curl -X POST http://localhost:8080/api/v1/wallet
 -H "Content-Type: application/json" -d '{"walletId":"11111111-1111-1111-1111-111111111111","operationType":"WITHDRAW","amount":200}'

Проверить баланс:
curl -X GET http://localhost:8080/api/v1/wallets/11111111-1111-1111-1111-111111111111
```

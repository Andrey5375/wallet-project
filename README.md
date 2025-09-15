# Wallet Service
Приложение на Go для работы с кошельками через REST API. Позволяет делать депозиты, снятия и получать баланс кошелька. Работает с PostgreSQL и запускается через Docker.

## Стек технологий
Golang, PostgreSQL, Docker & Docker Compose, REST API

## Запуск проекта
1. Склонировать репозиторий:  
`git clone <URL_вашего_репозитория>`  
`cd wallet-service`  

2. Создать файл `config.env` с переменными окружения:  
PORT=8080  
DATABASE_URL=postgres://user:password@db:5432/wallet?sslmode=disable  
POSTGRES_USER=user  
POSTGRES_PASSWORD=password  
POSTGRES_DB=wallet  

3. Поднять сервис и базу через Docker Compose:  
`docker-compose up --build`  

4. Приложение будет доступно на `http://localhost:8080`

## API
### POST /api/v1/wallet
Выполняет операцию DEPOSIT или WITHDRAW.  
**Request:**  
{
  "walletId": "UUID",
  "operationType": "DEPOSIT",
  "amount": 1000
}  
**Response:**  
{
  "status": "ok"
}

### GET /api/v1/wallets/{walletId}
Возвращает текущий баланс кошелька.  
**Response:**  
{
  "balance": 10000
}

## Тесты
Конкурентные тесты находятся в папке service_test.  
Запуск тестов:  
`go test -v ./service_test`

## Нагрузочное тестирование
Пример проверки 1000 запросов с hey:  
`hey -n 1000 -c 50 -m POST -H "Content-Type: application/json" -d '{"walletId":"11111111-1111-1111-1111-111111111111","operationType":"DEPOSIT","amount":10}' http://localhost:8080/api/v1/wallet`

## Особенности
Работа с конкурентными запросами (1000 RPS по одному кошельку), нет ошибок 50X при высокой нагрузке, переменные окружения читаются из config.env, полностью покрыто тестами на депозиты и снятия с учетом конкуренции.

# Auth Service
## Микросервис аутентификации

### Технологии: Go, PostgreSQL, Docker

#### Описание
Этот сервис реализует базовую аутентификацию пользователей с использованием JWT и refresh токенов.

Сервис полностью контейнеризирован и запускается одной командой.

Документация API доступна через Swagger.

## Быстрый старт
Требования
-Docker

-Docker Compose

-Go (для локальной разработки)

### Запуск
1.Склонируйте репозиторий:

git clone https://github.com/BroARD/MEDODS.git

cd /MEDODS

docker-compose -f docker-compose.yml up -d

2.Проверьте значения в .env (по умолчанию уже настроены).

3.Запустите сервис:

docker-compose -f docker-compose.yml up -d

4.Swagger-документация будет доступна по адресу:

http://localhost:8080/api/swagger/index.html

## Архитектура

-Go — реализация бизнес-логики и API

-PostgreSQL — хранение данных пользователей и refresh токенов (в виде bcrypt-хеша)

-Docker — контейнеризация приложения и базы данных

## API
### 1. Получение пары токенов

POST /api/tokens?guid=<user_guid>

Описание: Возвращает пару access/refresh токенов для пользователя с указанным GUID.

Пример запроса:

POST /api/tokens?guid=123e4567-e89b-12d3-a456-426614174000

Ответ:

json
{
  "access_token": "<jwt>",
  "refresh_token": "<base64>"
}

### 2. Обновление токенов

POST /api/refresh

Описание: Обновляет пару токенов. Требует действующую пару access/refresh токенов.

Требования:

User-Agent не должен меняться

При смене IP отправляется POST на WEBHOOK_URL

При ошибке — деавторизация пользователя

Пример запроса:

GET /api/v1/me

Authorization: Bearer <jwt>

json
{
  "refresh_token": "<base64>"
}

Ответ:

json
{
  "access_token": "<jwt>",
  "refresh_token": "<base64>"
}

###3. Получение GUID текущего пользователя

GET /api/tokens

Описание: Возвращает GUID текущего пользователя.

Требует авторизации (access token в заголовке Authorization).

Пример запроса:

GET /api/tokens
Authorization: Bearer <jwt>

Ответ:

json
{
  "user_id":"string"
}

### 4. Деавторизация пользователя

POST /api/logout

Описание: Деавторизует пользователя. После этого access и refresh токены становятся невалидными.

Пример запроса:
GET /api/tokens
Authorization: Bearer <jwt>

Ответ:
Logout complited!


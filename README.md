# Music Subscription Service

Backend REST API для музыкальной платформы (аналог Spotify/Яндекс Музыки), написанный на Go.

## Быстрый старт

### Требования
- [Docker](https://www.docker.com/) + Docker Compose
- Go 1.21+ (только для разработки)

### Запуск через Docker Compose

```bash
# Клонировать репозиторий
git clone <repo-url>
cd music-service

# Запустить все сервисы (backend + PostgreSQL + Redis)
docker-compose up --build

# Сервис доступен на http://localhost:8080
# Swagger UI: http://localhost:8080/swagger/index.html
```

### Запуск для разработки (без Docker)

```bash
# Установить зависимости
go mod download

# Настроить .env (скопировать пример)
cp .env.example .env

# Запустить PostgreSQL и Redis (через Docker)
docker-compose up postgres redis -d

# Применить миграции
make migrate-up

# Запустить приложение
make run
```

## Переменные окружения

| Переменная | По умолчанию | Описание |
|-----------|-------------|----------|
| `APP_PORT` | `8080` | Порт сервера |
| `ENV` | `development` | Среда (`development` / `production`) |
| `DB_HOST` | `localhost` | Хост PostgreSQL |
| `DB_PORT` | `5432` | Порт PostgreSQL |
| `DB_USER` | `music_user` | Пользователь БД |
| `DB_PASSWORD` | `music_password` | Пароль БД |
| `DB_NAME` | `music_service` | Имя базы данных |
| `REDIS_HOST` | `localhost` | Хост Redis |
| `REDIS_PORT` | `6379` | Порт Redis |
| `JWT_ACCESS_SECRET` | — | Секрет для access token (обязательно изменить!) |
| `JWT_REFRESH_SECRET` | — | Секрет для refresh token (обязательно изменить!) |
| `JWT_ACCESS_TTL` | `15m` | Время жизни access token |
| `JWT_REFRESH_TTL` | `720h` | Время жизни refresh token (30 дней) |
| `FREE_DAILY_PLAY_LIMIT` | `10` | Лимит прослушиваний в день для FREE |
| `FREE_PLAYLIST_LIMIT` | `3` | Лимит плейлистов для FREE |

## API Endpoints

### Auth
| Метод | Endpoint | Описание | Auth |
|-------|----------|----------|------|
| POST | `/api/v1/auth/register` | Регистрация | Нет |
| POST | `/api/v1/auth/login` | Авторизация | Нет |
| POST | `/api/v1/auth/refresh` | Обновление токена | Нет |
| POST | `/api/v1/auth/logout` | Выход | Нет |

### Users
| Метод | Endpoint | Описание | Auth |
|-------|----------|----------|------|
| GET | `/api/v1/users/me` | Профиль текущего пользователя | Да |

### Tracks
| Метод | Endpoint | Описание | Auth |
|-------|----------|----------|------|
| GET | `/api/v1/tracks` | Список треков (пагинация) | Да |
| GET | `/api/v1/tracks/{id}` | Трек по ID | Да |
| GET | `/api/v1/tracks/search?query=...` | Поиск треков | Да |
| POST | `/api/v1/tracks/{id}/play` | Прослушать трек | Да |

### Playlists
| Метод | Endpoint | Описание | Auth |
|-------|----------|----------|------|
| GET | `/api/v1/playlists` | Мои плейлисты | Да |
| POST | `/api/v1/playlists` | Создать плейлист | Да |
| GET | `/api/v1/playlists/{id}` | Плейлист по ID | Да |
| PUT | `/api/v1/playlists/{id}` | Обновить плейлист | Да |
| DELETE | `/api/v1/playlists/{id}` | Удалить плейлист | Да |
| POST | `/api/v1/playlists/{pid}/tracks/{tid}` | Добавить трек в плейлист | Да |
| DELETE | `/api/v1/playlists/{pid}/tracks/{tid}` | Убрать трек из плейлиста | Да |

### Favorites
| Метод | Endpoint | Описание | Auth |
|-------|----------|----------|------|
| GET | `/api/v1/favorites/tracks` | Избранные треки | Да |
| POST | `/api/v1/favorites/tracks/{id}` | Добавить в избранное | Да |
| DELETE | `/api/v1/favorites/tracks/{id}` | Удалить из избранного | Да |

### History
| Метод | Endpoint | Описание | Auth |
|-------|----------|----------|------|
| GET | `/api/v1/listening-history` | История прослушиваний | Да |

### Admin (требует роль ADMIN)
| Метод | Endpoint | Описание | Auth |
|-------|----------|----------|------|
| POST | `/api/v1/admin/tracks` | Добавить трек | Да (ADMIN) |
| PUT | `/api/v1/admin/tracks/{id}` | Обновить трек | Да (ADMIN) |
| DELETE | `/api/v1/admin/tracks/{id}` | Удалить трек | Да (ADMIN) |
| PATCH | `/api/v1/admin/users/{id}/subscription` | Изменить подписку | Да (ADMIN) |

## Авторизация

В заголовке запроса:
```
Authorization: Bearer <access_token>
```

В Swagger UI нажми **Authorize** и введи `Bearer <токен>`.

## Создание Admin пользователя

```bash
# Через psql в Docker
docker exec <postgres_container> psql -U music_user -d music_service

# Сгенерируй bcrypt хэш пароля, затем:
INSERT INTO users (email, username, password_hash, role, subscription_type)
VALUES ('admin@example.com', 'admin', '<bcrypt_hash>', 'ADMIN', 'PREMIUM');
```

Или используй `make create-admin EMAIL=admin@example.com PASSWORD=secret`.

## Makefile команды

```bash
make run          # Запустить приложение
make build        # Скомпилировать
make test         # Запустить тесты
make migrate-up   # Применить миграции
make migrate-down # Откатить миграции
make swag         # Сгенерировать Swagger docs
make docker-up    # docker-compose up --build
make docker-down  # docker-compose down
make lint         # Линтер
```

## Структура базы данных

- `users` — пользователи
- `refresh_tokens` — токены обновления
- `artists`, `albums`, `genres` — музыкальный каталог
- `tracks` — треки
- `playlists`, `playlist_tracks` — плейлисты
- `favorites` — избранное
- `listening_history` — история прослушиваний

## Технологии

| Компонент | Технология |
|-----------|-----------|
| Язык | Go 1.21 |
| Веб-фреймворк | [chi](https://github.com/go-chi/chi) |
| База данных | PostgreSQL 16 |
| ORM/DB | [sqlx](https://github.com/jmoiron/sqlx) |
| Кэш | Redis 7 |
| Авторизация | JWT (golang-jwt) |
| Логирование | [zap](https://github.com/uber-go/zap) |
| Документация | Swagger (swaggo) |
| Миграции | [goose](https://github.com/pressly/goose) |
| Контейнеризация | Docker Compose |

## Swagger

Swagger UI доступен после запуска по адресу:
```
http://localhost:8080/swagger/index.html
```

## Тесты

```bash
go test ./... -v
```

## Архитектура

```
cmd/app/main.go              # Точка входа
internal/
  core/                      # Переиспользуемые компоненты
    config/                  # Конфигурация
    logger/                  # Логирование (zap)
    middleware/              # HTTP middleware
    postgres/                # Подключение к БД
    redis/                   # Подключение к Redis
    response/                # HTTP ответы
    errors/                  # Общие ошибки
  features/                  # Бизнес-логика по фичам
    auth/                    # Аутентификация
    tracks/                  # Треки
    playlists/               # Плейлисты
    favorites/               # Избранное
    history/                 # История
    users/                   # Пользователи
    subscriptions/           # Подписки
migrations/                  # SQL миграции
docs/                        # Swagger документация
frontend/                    # Веб-интерфейс
```

## Деплой на Fly.io и CI/CD

Проект настроен для автоматического деплоя на платформу **Fly.io** при коммите в ветку `main` в GitLab.

### Локальный запуск деплоя

1. Установите CLI-инструмент `flyctl`:
   ```bash
   brew install flyctl
   ```
2. Авторизуйтесь под своим аккаунтом:
   ```bash
   flyctl auth login
   ```
3. Инициализируйте проект в Fly (если приложение еще не создано):
   ```bash
   flyctl launch
   ```
4. Создайте базу данных Postgres и Redis в Fly.io и свяжите их с приложением, либо настройте переменные окружения вручную:
   ```bash
   flyctl secrets set DB_HOST=... DB_USER=... DB_PASSWORD=... DB_NAME=...
   flyctl secrets set REDIS_HOST=... REDIS_PORT=... REDIS_PASSWORD=...
   flyctl secrets set JWT_ACCESS_SECRET=... JWT_REFRESH_SECRET=...
   ```
5. Запустите деплой вручную:
   ```bash
   flyctl deploy
   ```

### GitLab CI/CD

В проекте настроен пайплайн в файле `.gitlab-ci.yml`, состоящий из двух этапов:
1.  **test:** Запуск всех юнит-тестов (интеграционные тесты БД пропускаются в режиме `-short`).
2.  **deploy:** Автоматический деплой на Fly.io при слиянии/пуше в ветку `main`.

Для работы автоматического деплоя необходимо добавить токен авторизации в GitLab CI/CD:
1. Зайдите в GitLab: **Settings** → **CI/CD** → **Variables**.
2. Добавьте переменную `FLY_API_TOKEN`. В качестве значения укажите токен доступа, полученный через команду `flyctl auth token` или в панели управления Fly.io.


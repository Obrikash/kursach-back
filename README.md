# REST API for swimming pools (CRM-like)

## Installation and Start

Перед запуском сервера, нужно создать .envrc файл, куда поместить переменные окружения требуемые в Makefile.
(SWIMMING_POOL_DSN, JWT_SECRET)

```bash
$ docker-compose up -d
$ make run/api
```

Первая команда запускает **PostgreSQL** в контейнере, вторая запускает само приложение. Позже в docker-compose будет добавлено и само приложение

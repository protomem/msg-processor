# Msg Processor

## Описание

Сервис для обработки сообщений.

- Сылка на работующий проект: [тык](http://80.90.184.101/swagger)
- Доступ к админ панели кафки: [тык](http://80.90.184.101:9000)
- Доступ к админ панели БД: [тык](http://80.90.184.101:15432)
  - Логин и пароль: `admin@pgadmin.com` `123456789`
  - Имя хоста БД: `pg`
  - Логин и пароль к БД (по умолчанию): `admin` `123456789`

## Используемые технологии

- Go 1.22
- Postgres 14
- Kafka

## Конфигурация

Основная конфигурация хранится в переменных окружения.
Флаг (опциональный) `-cfg` указывает путь к конфигурационному файлу(например, `.env`).

### Переменные окружения

- `STORE_DSN` - строка подключения к базе данных
- `STORE_PING` - флаг проверки подключения к базе данных
- `QUEUE_ADDRS` - адреса кафки
- `QUEUE_TOPIC` - топик кафки
- `LISTEN_ADDR` - адрес для прослушивания сервера
- `BASE_URL` - адрес для доступа к API
- `READ_PROC_MSGS_INTERVAL` - интервал опроса кафки
- `READ_PROC_MSGS_TIMEOUT` - время ожидания сообщения от кафки

## Запуск

### Клонирование

```sh
git clone https://github.com/protomem/msg-processor.git && cd msg-processor
```

### Docker

```sh
docker compose up -d
# или
make run-docker
```

## Миграции

```sh
DB_DSN="<db_dsn>" make migrations/up # запустить миграции

DB_DSN="<db_dsn>" make migrations/down # откатить миграции

DB_DSN="<db_dsn>" make migrations/new name=<name> # создать новую миграцию

DB_DSN="<db_dsn>" make migrations/goto version=<version> # перейти на версию <version> миграции

DB_DSN="<db_dsn>" make migrations/force version=<version> # применить миграцию версии <version>
```

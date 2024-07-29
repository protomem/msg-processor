# Msg Processor

## Описание

Тестовое задание. Сервис для обработки сообщений.

- Сылка на работующий проект: [тык](http://80.90.184.101/swagger)

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
- `STORE_MIGRATE` - автоматическая миграция (по-умолчанию `true`)
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
```
или
```sh
make run-docker
```

- Не забудьте, если вы запускаете приложение не для тестирования, вам нужно указать в `BASE_URL` IP-адрес вашего сервера.

## Миграции

```sh
DB_DSN="<db_dsn>" make migrations/up # запустить миграции

DB_DSN="<db_dsn>" make migrations/down # откатить миграции

DB_DSN="<db_dsn>" make migrations/new name=<name> # создать новую миграцию

DB_DSN="<db_dsn>" make migrations/goto version=<version> # перейти на версию <version> миграции

DB_DSN="<db_dsn>" make migrations/force version=<version> # применить миграцию версии <version>
```

## FAQ

### Что, если Kakfka не читает сообщения?

```sh
docker compose down && docker compose up -d
```
  

# ACOMM - Сервис комментариев

Сервис для хранения и работы с комментариями к новостям. Использует SQLite для хранения данных.

## API Сервиса

Сервис работает на порту 8082 и предоставляет следующие API:

### Получение комментариев для новости

```
GET /api/comm_news?id={newsId}
```

Где `{newsId}` - идентификатор новости.

#### Ответ

```json
[
  {
    "id": 1,
    "news_id": 1,
    "text": "Это первый комментарий"
  },
  {
    "id": 2,
    "news_id": 1,
    "text": "Это второй комментарий"
  }
]
```

Статус 200 OK при успешном запросе.

### Добавление комментария к новости

```
POST /api/comm_add_news?id={newsId}
```

Где `{newsId}` - идентификатор новости.

#### Запрос

```json
{
  "text": "Текст комментария"
}
```

#### Ответ

```json
{
  "id": 3
}
```

Где `id` - идентификатор нового комментария.

Статус 201 Created при успешном создании.

## Запуск сервиса

```bash
# Сборка
go build -o acomm ./cmd/server

# Запуск
./acomm
```

База данных автоматически создается при первом запуске по пути `./db/comm.db`. 
[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/xR-tWBKa)

Микросервис для работы с переговорками (создание переговорок, просмотр списка переговорок, создание расписания для
переговорок, просмотр доступных слотов для брони, отмена брони, просмотр списка своих броней или же всех броней с
пагинацией

## Стек технологий

| Компонент           | Технология                           |
|---------------------|--------------------------------------|
| **Язык**            | Go 1.25                              |
| **Фреймворк**       | Gin                                  |
| **База данных**     | PostgreSQL 15                        |
| **ORM / SQL**       | GORM + lib/pq                        |
| **Миграции**        | golang-migrate (встроены в бинарник) |
| **Контейнеризация** | Docker & Docker Compose              |
| **Тестирование**    | testify, testcontainers-go           |
| **Build**           | Makefile                             |

Сервис был написан, стараясь придерживаться (не везде получилось реализовать задуманное) принципа "Чистая архитектура", 
с целью легкого расширения функционала и простотой реализации юнит-тестирования. Также был реализован Graceful Shutdown 
для корректного завершения работы сервиса


## Архитектурные решения

### 1. Генерация слотов: при запросе

**Проблема:** Каким образом необходимо генерировать доступные слоты на основе расписания соответствующей переговорки?

**Решение:** Было решено генерировать слоты **«по требованию»** при первом запросе `/slots/list` на конкретную дату.
Однако было бы лучше использовать гибридный подход (генерировать по требованию и хранить слоты только в кеше, 
без таблицы в БД, т.к. нам необходимо было бы только `(7 дней*1000 слотов/день)=7000 слотов` в кеше).
По поводу скользящего окна - это дополнительное усложнение при реализации)

### 2. Одновременное(конкурентное) бронирование слота

**Проблема:** Два пользователя могут одновременно попытаться забронировать один и тот же слот.

**Решение:** Использование транзакции БД с блокировкой строки (`SELECT ... FOR UPDATE`).

# Запуск

Для запуска сервиса необходимо предварительно заполнить .env файл

Запустить сервис можно с помощью команды `make up`

Для запуска тестов с покрытием выполнить команду `make test-coverage`. Для получения отчета в формате HTML выполнить
команду `test-coverage` для получения отчёта в html формате


## Примеры запросов

### Примитивная авторизация (dummyLogin)

Авторизация:

```curl
curl --location 'http://localhost:8080/dummyLogin' \
--header 'Content-Type: application/json' \
--data '{
    "role": "admin"
}'
```

Пример ответа:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAxIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzc0NTM2MDg1LCJpYXQiOjE3NzQ0NDk2ODV9.sWdvUQrti-i9calpmqypdxEw5NcDxBmNLpVi3EYsjTw"
}
```

### Создание переговорки

Создание переговорки админом:

```curl
curl --location 'http://localhost:8080/rooms/create' \
--header 'Content-Type: text/plain' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAxIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzc0NTM2MDg1LCJpYXQiOjE3NzQ0NDk2ODV9.sWdvUQrti-i9calpmqypdxEw5NcDxBmNLpVi3EYsjTw' \
--data '{
    "name": "Room #5",
    "capacity": 15
}'
```

Пример ответа:

```json
{
  "room": {
    "id": "f8c2eadf-af8b-4e46-948d-16a171fa4743",
    "name": "Room #5",
    "capacity": 15,
    "createdAt": "2026-03-25T14:54:38.942316554Z"
  }
}
```

### Список переговорок

Список всех переговорок:

```curl
curl --location 'http://localhost:8080/rooms/list' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAyIiwicm9sZSI6InVzZXIiLCJleHAiOjE3NzQ1MzYwODksImlhdCI6MTc3NDQ0OTY4OX0.pRhSGCv_OO_jW6mUh3ASl7f49gYk6KEXWgAmjYbq9Xk'
```

Пример ответа:

```json
{
  "rooms": [
    {
      "id": "f8c2eadf-af8b-4e46-948d-16a171fa4743",
      "name": "Room #5",
      "capacity": 15,
      "createdAt": "2026-03-25T14:54:38.942316Z"
    }
  ]
}
```

### Создание расписания для переговорки

Новое расписание переговорки создает админ:

```curl
curl --location 'http://localhost:8080/rooms/f8c2eadf-af8b-4e46-948d-16a171fa4743/schedule/create' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAxIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzc0NTM2MDg1LCJpYXQiOjE3NzQ0NDk2ODV9.sWdvUQrti-i9calpmqypdxEw5NcDxBmNLpVi3EYsjTw' \
--data '{
    "daysOfWeek": [
        1,
        2,
        3,
        4,
        5
    ],
    "startTime": "9:30",
    "endTime": "13:00"
}'
```

Пример ответа:

```json
{
  "schedule": {
    "id": "db3a7ab6-2809-4707-b2df-94a7731dc78b",
    "roomId": "f8c2eadf-af8b-4e46-948d-16a171fa4743",
    "daysOfWeek": [
      1,
      2,
      3,
      4,
      5
    ],
    "startTime": "9:30",
    "endTime": "13:00"
  }
}
```

### Список доступных слотов для переговорки с некоторым id

Список слотов, которые соответствуют расписанию переговорки с некоторым id:

```curl
curl --location 'http://localhost:8080/rooms/f8c2eadf-af8b-4e46-948d-16a171fa4743/slots/list?date=2026-03-26' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAxIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzc0NTM2MDg1LCJpYXQiOjE3NzQ0NDk2ODV9.sWdvUQrti-i9calpmqypdxEw5NcDxBmNLpVi3EYsjTw'
```

Пример ответа:

```json
{
  "slots": [
    {
      "id": "81038413-38ce-4e55-938c-8c01e5f03a72",
      "roomId": "f8c2eadf-af8b-4e46-948d-16a171fa4743",
      "start": "2026-03-26T09:30:00Z",
      "end": "2026-03-26T10:00:00Z"
    },
    {
      "id": "6113e38b-dca3-4bb8-ac48-4e16e97ce29b",
      "roomId": "f8c2eadf-af8b-4e46-948d-16a171fa4743",
      "start": "2026-03-26T10:00:00Z",
      "end": "2026-03-26T10:30:00Z"
    }
  ]
}
```

### Бронирование переговорки

Бронируем переговорку по доступному слоту:

```curl
curl --location 'http://localhost:8080/bookings/create' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAyIiwicm9sZSI6InVzZXIiLCJleHAiOjE3NzQ1MzYwODksImlhdCI6MTc3NDQ0OTY4OX0.pRhSGCv_OO_jW6mUh3ASl7f49gYk6KEXWgAmjYbq9Xk' \
--data '{
    "slotId": "6113e38b-dca3-4bb8-ac48-4e16e97ce29b",
    "createConferenceLink": true
}'
```

Пример ответа:

```json
{
  "booking": {
    "id": "b9a90e36-1155-42e2-9925-03ac82607226",
    "slotId": "6113e38b-dca3-4bb8-ac48-4e16e97ce29b",
    "userId": "00000000-0000-0000-0000-000000000002",
    "status": "active",
    "conferenceLink": "https://conference.com/meeting/b9a90e36-1155-42e2-9925-03ac82607226",
    "createdAt": "2026-03-25T14:55:34.44710359Z"
  }
}
```

### Отмена брони

Отменяем бронь по его id

```curl
curl --location --request POST 'http://localhost:8080/bookings/d5c05f32-a2bf-401e-96e9-f5907c16b5a7/cancel' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAyIiwicm9sZSI6InVzZXIiLCJleHAiOjE3NzQ1MzYwODksImlhdCI6MTc3NDQ0OTY4OX0.pRhSGCv_OO_jW6mUh3ASl7f49gYk6KEXWgAmjYbq9Xk'
```

Пример ответа:

```json
{
  "booking": {
    "id": "d5c05f32-a2bf-401e-96e9-f5907c16b5a7",
    "slotId": "81038413-38ce-4e55-938c-8c01e5f03a72",
    "userId": "00000000-0000-0000-0000-000000000002",
    "status": "cancelled",
    "createdAt": "2026-03-25T14:55:17.491433Z"
  }
}
```

### Список всех броней пользователя

Пользователь получает весь список своих броней:

```curl
curl --location 'http://localhost:8080/bookings/my' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAyIiwicm9sZSI6InVzZXIiLCJleHAiOjE3NzQ1MzYwODksImlhdCI6MTc3NDQ0OTY4OX0.pRhSGCv_OO_jW6mUh3ASl7f49gYk6KEXWgAmjYbq9Xk'
```

Пример ответа:

```json
{
  "bookings": [
    {
      "id": "b9a90e36-1155-42e2-9925-03ac82607226",
      "slotId": "6113e38b-dca3-4bb8-ac48-4e16e97ce29b",
      "userId": "00000000-0000-0000-0000-000000000002",
      "status": "active",
      "conferenceLink": "https://conference.com/meeting/b9a90e36-1155-42e2-9925-03ac82607226",
      "createdAt": "2026-03-25T14:55:34.447103Z"
    },
    {
      "id": "d5c05f32-a2bf-401e-96e9-f5907c16b5a7",
      "slotId": "81038413-38ce-4e55-938c-8c01e5f03a72",
      "userId": "00000000-0000-0000-0000-000000000002",
      "status": "cancelled",
      "createdAt": "2026-03-25T14:55:17.491433Z"
    }
  ]
}
```

### Список броней с пагинацией

Админ выводит список броней с пагинацией (размер страницы и номер страницы):

```curl
curl --location 'http://localhost:8080/bookings/list?pageSize=2&page=1' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAxIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzc0NTM2MDg1LCJpYXQiOjE3NzQ0NDk2ODV9.sWdvUQrti-i9calpmqypdxEw5NcDxBmNLpVi3EYsjTw'
```

Пример ответа:

```json
{
  "bookings": [
    {
      "id": "b9a90e36-1155-42e2-9925-03ac82607226",
      "slotId": "6113e38b-dca3-4bb8-ac48-4e16e97ce29b",
      "userId": "00000000-0000-0000-0000-000000000002",
      "status": "active",
      "conferenceLink": "https://conference.com/meeting/b9a90e36-1155-42e2-9925-03ac82607226",
      "createdAt": "2026-03-25T14:55:34.447103Z"
    },
    {
      "id": "d5c05f32-a2bf-401e-96e9-f5907c16b5a7",
      "slotId": "81038413-38ce-4e55-938c-8c01e5f03a72",
      "userId": "00000000-0000-0000-0000-000000000002",
      "status": "cancelled",
      "createdAt": "2026-03-25T14:55:17.491433Z"
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 2,
    "total": 2
  }
}
```
# Subscriptions service

REST-сервис для агрегации данных об онлайн-подписках пользователей.

## Запуск

Для запуска требуется .env. Пример:
```
DATABASE_URL=postgres://postgres:postgres@db:5432/subscriptions?sslmode=disable
PORT=8080
MIGRATIONS_PATH=./internal/db/migrations/
```

Написал небольшую утилиту для исполнения миграций (при запуске НЕ в docker). В docker всё выполняется само и MIGRATIONS_PATH не нужен.

При запуске сервиса также запускается swagger документация, расположенная по http://localhost:PORT/swagger/index.html. 

## Endpoints
[GET] /subscriptions - Список подписок
[POST] /subscriptions - Добавление подписки
[GET] /subscriptions/{id} - Получение подписки по id (SERIAL PRIMARY KEY)
[GET] /subscriptions/user/{user_id} - Получение подписок по user id (UUID)
[PUT] /subscriptions/{id} - Обновление данных о подписке по id (SERIAL PRIMARY KEY)
[DELETE] /subscriptions/{id} - Удаление подписки по id (SERIAL PRIMARY KEY)

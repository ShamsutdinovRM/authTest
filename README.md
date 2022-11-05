# Тестовое задание

## Регистрация и аутентификация.

**Задача:**

Необходимо написать микросервис авторизации на go, в котором реализованы следующие функции:
- создание аккаунта
- прохождение авторизации
- выход из аккаунта

При создании аккаунта необходимо записывать данные для авторизации в бд

При прохождении авторизации необходимо создавать токен и сохранять его куда-либо

При выходе из аккаунта необходимо удалять сохранённый токен пользователя

## Как использовать ##

**Создание и запуск докер образа с миграциями**
```sh
docker-compose build
docker-compose up
Создать таблицу на основе миграций
```

**Примеры запросов к сервису**

Сервис тестировался с использованием Postman.

**/signup**

Запрос:
```sh
{
    "username": "Protegos",
    "password": "123"
}
```
Ответ:
```sh
"Protegos"
```
Запрос с ошибкой:
```sh
{
    "username": "Protegos",
    "password": 123
}
```
Ответ:
```sh
{
    "text": "Invalid field"
}
```

**/login**

Запрос:
```sh
{
    "username": "Protegos",
    "password": "123"
}
```
Ответ:
```sh
{
    "username": "Protegos",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Njc2NTA1ODUsImlhdCI6MTY2NzY1MDI4NSwidXNlciI6IlByb3RlZ29zIn0.fUQCcEmXBpCo8dfqaO9-VgEZbcXHiQ_gJuvqNY7rN7s"
}
```
Запрос с ошибкой:
```sh
{
    "username": "Patrego",
    "password": "123"
}
```
Ответ:
```sh
{
    "text": "error, user not found: sql: no rows in result set"
}
```

**/auth/hello**

Чтобы получить доступ к данной странице, необходимо вставить в поле Bearer Token токен, полученный при Login, на вкладке Authorization

Запрос:
```sh
{
    "username": "Protegos",
    "password": "123"
}
```
Ответ:
```sh
"Hello From Wrong Side Of Heaven Protegos"
```

**/auth/logout**

Чтобы получить доступ к данной странице, необходимо вставить в поле Bearer Token токен, полученный при Login, на вкладке Authorization
Если вы не авторизованы, то и выйти у вас не выйдет

Запрос:
```sh
{
    "username": "Protegos"
}
```
Ответ:
```sh
"Goodbye"
```

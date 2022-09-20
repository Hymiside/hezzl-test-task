# Сервис для тестового задания в [Hezzl.com](https://hezzl.com/)

## Как запустить сервис
1. Склонируйте репозиторий
2. Установите зависимости
```
~$ go mod vendor
```
3. Установите переменные окружения в файле `.env`
4. Накатите миграцию `schema/postgres/0001.up.sql`
5. Запустите файл `cmd/main.go`

## Методы
1. Create Item - создает новый айтем зависимый, обновляет Redis и возвращает созданный айтем
```json
TYPE: POST
URL: /item/create?campaignId=int

Payload: 
{
  "name": "string", // обязательное поле
  "description": "string" // обязательное поле
}
```
2. Update item - обновляет данные айтема, обновляет Redis и возвращает обновленные данные айтема
```
TYPE: PATCH
URL: /item/update?id=int&campaignId=int

Payload: {
  "name": "string", // обязательное поле
  "description": "string" // необязательное поле
}
```
3. Delete Item - удаляет item, обновляет Redis и возвращает данные об удаленном айтеме
```
TYPE: DELETE
URL: /item/remove?id=int&campaignId=int

Payload: {}
```
4. List Item - возвращает все неудаленные айтемы. Если в Redis есть данные, то возвращает оттуда, если нет, то достает из Postgres, кэширует в Redis на 1 минуту и возвращает
```
TYPE: GET
URL: /item/list

Payload: {}
```
## Логи
Логи пишутся пачками по 24 штуки в Postgres через Nats

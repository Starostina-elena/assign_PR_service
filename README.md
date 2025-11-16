# **Тестовое задание для стажёра Авито Backend (осенняя волна 2025)**

## **Сервис назначения ревьюеров для Pull Request’ов**
## Выполнила Старостина Елена

### Запуск
#### Перед запуском заполните .env (hint: достаточно переименовать .env.example в .env)
```make up```

### Остановка
```make down```

### Завершение работы
```make clean```

### Примеры запросов 
#### Создание юзера 
```
curl -i -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"id":"u1","name":"Alice","is_active":true}'
```
#### Получение юзера 
```
curl http://localhost:8080/users/u1      
```
#### Активация/деактивация юзера
```
curl -i -X PUT http://localhost:8080/users/setIsActive -H "Content-Type: application/json" -d '{"user_id":"u1","is_active":false}'
```

#### Получение открытых PR, где пользователь назначен ревьюером 
```
curl http://localhost:8080/users/u1/getReview
```
#### Создание команды
```
curl -i -X POST 'http://localhost:8080/team/add' -H 'Content-Type: application/json' -d '{
    "team_name": "payments",
    "members": [
      {"user_id": "u1", "username": "Alice", "is_active": true},
      {"user_id": "u2", "username": "Bob",   "is_active": true}
    ]
  }'
```
#### Получение команды
```
curl http://localhost:8080/team/1/get
```
#### Создание PR
```
curl -i -X POST 'http://localhost:8080/pullRequest/create' -H 'Content-Type: application/json' -d '{
    "pull_request_name": "Add search",
    "author_id": "u1"
  }'
```
#### Reassign PR
```
curl -sS -X POST http://localhost:8080/pullRequest/reassign -H "Content-Type: application/json" -d '{"pull_request_id":1,"old_user_id":"u2"}'
```
#### Merge PR
```
curl -i -X PUT 'http://localhost:8080/pullRequest/merge' -H 'Content-Type: application/json' -d '{"pull_request_id": 1}'
```

## Допущения, принятые в ходе работы:
- Для users/getReview не был задан конкретный путь с использованием UserId. Так как описанием API также не было указано на необходимость использования тела запроса, было принято решение поместить userId в url запроса после users/
- Также указанный getReview возвращает не все PR, а только открытые, потому что это показалось мне логичным. При необходимости фиксится снятием условия на is_opened = TRUE в internal/storage/user_repo.GetPullRequestsAssigned
- У users также есть дополнительные эндпоинты (создание юзера, получение юзера по id, установка команды юзеру, удаление юзера из команды) для удобства тестирования 
- В базе данных в целях оптимизации для PR используется не enum OPEN|MERGED, а булево is_opened. Для удобства пользователя is_opened заменяется на OPEN|MERGED при отдаче пользователю. 
- Ближе к концу выполнения работы, внимательно изучая openapi.yml, с удивлением обнаружила, что для юзеров и PR в качестве id используются строки. Для юзеров исправила, для PR не успела Т_Т 
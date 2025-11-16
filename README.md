# **Тестовое задание для стажёра Авито Backend (осенняя волна 2025)**

## **Сервис назначения ревьюеров для Pull Request’ов**


### Запуск
```make up```

### Остановка
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

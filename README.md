

# Веб-сервис для вычисления арифметических выражений

## Описание
Проект для яндекс лицея, вычисление математических выражений.

## Запуск
1. Выполнить команду go mod tidy
2. Запустить сервер.
#### `go run ./cmd/main.go `
3. Сервис будет доступен по адресу: [http://localhost:8080/api/v1/calculate](http://localhost:8080/api/v1/calculate). 
С одним эндпоинтом, котороый принимает JSON с математическим выражением. `POST /api/v1/calculate`

## Пример запроса с использованием curl
`curl -X POST http://localhost:8080/api/v1/calculate -H "Content-Type: application/json" -d "{\"expression\": \"25/(6-5)+(1+2)*15\"}"`
### Примеры ответов

1. **Успешный запрос**:
Ответ:

```json
{"result": "6"}
```

2. **Ошибка: некорректное выражение**:

```json
{"error":"Expression is not valid"}
```

В зависимости от выражения будет разный ответ, примеры можно посмотреть в тестах.
## Команды для тестирования
1. Тестирование калькулятора
```
cd .\pkg\calculation\ 
go test
```
2. Тестирование сервера
```
cd .\internal\application\
go test
```
## Установить нужный вам порт
```powershell
$env:PORT=9090
go run ./cmd/main.go
```

## Почта для связи
LarkovAleksandr005@yandex.ru

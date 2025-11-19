# JSONPlaceholder API

REST API для работы с постами через JSONPlaceholder.

## Описание

Микросервис для JSONPlaceholder API с валидацией, логированием и обработкой ошибок.

## Клонирование репозитория

```bash
git clone https://github.com/BochkaDeyalo/jsonplaceholder.api.git
```

## Запуск

```bash
go run main.go
```

Сервер запустится на http://localhost:8080

## Сборка

```bash
go build main.go
```

## Swagger

Документация доступна по адресу:

```
http://localhost:8080/swagger/index.html
```

Генерация документации после изменения аннотаций:

```bash
swag init
```

## Тестирование

Запуск тестов контроллера:
```bash
go test ./controller/ -v
```

Проверка покрытия:
```bash
go test ./controller/ -cover
```

Детальный отчет:
```bash
go test ./controller/ -coverprofile=coverage.out
```

## Валидация

Правила валидации:
- userId - обязательное, минимум 1
- title - обязательное, минимум 1 символ
- body - обязательное, минимум 1 символ

## Обработка ошибок

Формат ошибок:
```json
{
  "code": 400,
  "message": "Validation failed",
  "details": "подробное описание ошибки"
}
```

Коды ошибок:
- 400 - некорректный JSON или ошибка валидации
- 502 - ошибка внешнего сервиса

## Логирование

Уровни логирования:
- DEBUG - детальная информация
- INFO - общая информация
- WARN - предупреждения
- ERROR - ошибки

Настройка через переменную LOG_LEVEL в .env файле. 

## Архитектура

Слои приложения:
1. Controller - обработка HTTP запросов
2. Service - слой бизнес-логики (работа с внешним API)
3. Model - модели данных
4. Error - типизированные ошибки
5. Logger - зависимость логгер

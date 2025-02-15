# MerchShop

# Стек
- Gorm — для работы с базой данных (ORM).
- PostgreSQL — для работы с реляционными базами данных.
- SQLite — для легковесной базы данных, используемой в тестах.
- gomock — для создания моков в юнит-тестах.
- testify — для написания тестов с удобными утверждениями (asserts).
- Gin — для создания API.
- Jwt — для генерации и валидации JWT токенов.
- go-sqlmock — для мокирования SQL запросов в тестах.

 
## Инструкция по запуску

Выполнить команду в дериктории c файлом docker-compose.yaml

`docker-compose up --build -d`

Для проверки работоспособности можно тестировать с помощью Swagger Editor, используя [API](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-winter-2025/schema.json) из условий задания.

## Проверка покрытия текстами

`go test ./... -coverprofile=coverage | Out-Null`

`go tool cover -func=coverage | Select-String "total"`

![image](https://github.com/user-attachments/assets/f9a72ead-417a-4382-ae6c-a2ec4ad2ab9f)

## Итоги

- Изучил Unit и интеграционное тестирование
- Практиковался в использовании чистой архитектуры в проекте, который реализует API
- Использовал SwaggerEditor, чтобы точно следовать техническому заданию.
- Вспомнил работу с Docker
- Изучил механизм Graceful Shutdown





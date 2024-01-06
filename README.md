## Инструкция по запуску
### 1. Сборка Docker образа
    docker build -t todo-service .

### 2. Запуск Сервисов
**Примечение:** сначало проект выдаст ошибку потому, что еще не было миграций.

    docker-compose up --build

### 3. Запуск Миграций
**Примечение:** миграции надо запускать пока включен контейнер.

    docker exec todo_api go run /app/migrate/migrate.go

### 4. Запуск Тестов
    docker exec todo_api go test /app/tests

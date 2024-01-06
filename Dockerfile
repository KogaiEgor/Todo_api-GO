# Использование базового образа с Go
FROM golang:latest

# Установка рабочей директории в контейнере
WORKDIR /app

# Копирование go модулей и их установка
COPY go.mod ./
COPY go.sum ./
COPY .env ./
RUN go mod download 

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN go build -o /todo-service

# Запуск приложения
CMD [ "/todo-service" ]
# Указываем версию golang для большей предсказуемости и стабильности сборки
FROM golang:1.22 AS compiling_stage

# Создаем и устанавливаем рабочую директорию
WORKDIR /news_agregator

# Копируем файлы go.mod и go.sum перед кодом приложения для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код приложения
COPY . .

# Сборка приложения
RUN go build -o news_agregator ./cmd/main.go

# Используем минимальный образ для запуска контейнера
FROM alpine:3.17

# Устанавливаем метаданные
LABEL version="1.0.0"

# Устанавливаем рабочую директорию
WORKDIR /root/

# Установите libc6-compat для совместимости с некоторыми Go бинарниками
RUN apk add --no-cache libc6-compat

# Копируем бинарный файл из предыдущего этапа
COPY --from=compiling_stage /news_agregator .

# Указываем точку входа
ENTRYPOINT ["./news_agregator"]

# Открываем порт 8080
EXPOSE 8080

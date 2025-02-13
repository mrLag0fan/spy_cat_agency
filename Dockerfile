# Використовуємо офіційний образ Golang для побудови програми
FROM golang:1.23 AS builder

# Встановлюємо робочу директорію в контейнері
WORKDIR /app

# Копіюємо go.mod та go.sum для кешування залежностей
COPY go.mod go.sum ./
RUN go mod tidy

# Копіюємо всю директорію проекту в контейнер
COPY . .

RUN apt-get update && apt-get install -y --no-install-recommends netcat-openbsd && rm -rf /var/lib/apt/lists/*

# Створюємо виконуваний файл (зазначаємо, що основний файл у папці cmd)
WORKDIR /app/cmd

# Будуємо програму, вказуючи на основний файл
RUN go build -o main .

COPY scripts/wait-for-db.sh /usr/local/bin/wait-for-db.sh
RUN chmod +x /usr/local/bin/wait-for-db.sh

# Запускаємо програму в контейнері
CMD ["/usr/local/bin/wait-for-db.sh", "db", "5432", "/app/cmd/main"]

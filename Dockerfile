FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

# ========= Runtime stage =========
FROM python:3.11-slim

# Install system dependencies for psycopg2 and pdfplumber
RUN apt-get update && apt-get install -y \
    libpq-dev gcc musl-dev python3-dev build-essential \
 && pip install --no-cache-dir pdfplumber psycopg2-binary python-dotenv \
 && apt-get remove -y gcc build-essential python3-dev \
 && apt-get autoremove -y && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/main .

COPY web-scraper/ ./web-scraper/

EXPOSE 8080

CMD ["./main"]


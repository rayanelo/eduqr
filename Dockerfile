# Étape 1 : Compilation avec Go 1.23
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copier les fichiers de dépendances
COPY go.mod go.sum ./
RUN go mod download

# Copier le reste du code
COPY . .

# Compiler le binaire
RUN go build -o main ./cmd/server

# Étape 2 : Image finale légère
FROM alpine:latest

# Installer curl pour les health checks
RUN apk add --no-cache curl

WORKDIR /root/

# Copier le binaire compilé
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"] 